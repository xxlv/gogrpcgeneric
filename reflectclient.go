package gogrpcgeneric

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/xxlv/gogrpcgeneric/dialer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

var exit = os.Exit

type diaConfig struct {
	connectTimeout float64
	keepaliveTime  float64
	maxMsgSz       int
	ua             string
	isUnixSocket   bool
	target         string
}

var base64Codecs = []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding}

func decode(val string) (string, error) {
	var firstErr error
	var b []byte
	// we are lenient and can accept any of the flavors of base64 encoding
	for _, d := range base64Codecs {
		var err error
		b, err = d.DecodeString(val)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return string(b), nil
	}
	return "", firstErr
}

func MetadataFromHeaders(headers []string) metadata.MD {
	md := make(metadata.MD)
	for _, part := range headers {
		if part != "" {
			pieces := strings.SplitN(part, ":", 2)
			if len(pieces) == 1 {
				pieces = append(pieces, "") // if no value was specified, just make it "" (maybe the header value doesn't matter)
			}
			headerName := strings.ToLower(strings.TrimSpace(pieces[0]))
			val := strings.TrimSpace(pieces[1])
			if strings.HasSuffix(headerName, "-bin") {
				if v, err := decode(val); err == nil {
					val = v
				}
			}
			md[headerName] = append(md[headerName], val)
		}
	}
	return md
}

type compositeSource struct {
	reflection DescriptorSource
}

func (cs compositeSource) ListServices() ([]string, error) {
	return cs.reflection.ListServices()
}

func (cs compositeSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	d, err := cs.reflection.FindSymbol(fullyQualifiedName)
	if err == nil {
		return d, nil
	}

	return nil, nil
}

func (cs compositeSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	exts, err := cs.reflection.AllExtensionsForType(typeName)
	if err != nil {
		// On error fall back to file source
		panic("error ")
	}
	// Track the tag numbers from the reflection source
	tags := make(map[int32]bool)
	for _, ext := range exts {
		tags[ext.GetNumber()] = true
	}
	return exts, nil
}

func GetDescSource(ctx context.Context, refClient *grpcreflect.Client) (compositeSource, error) {
	reflSource := DescriptorSourceFromServer(ctx, refClient)
	// 获取描述
	descSource := compositeSource{reflSource}
	return descSource, nil
}

// 获取反射的client
func doGetReflectClient(ctx context.Context, cc *grpc.ClientConn) *grpcreflect.Client {
	md := MetadataFromHeaders([]string{"__ref_version:1.0"})

	refCtx := metadata.NewOutgoingContext(ctx, md)
	refClient := grpcreflect.NewClientV1Alpha(refCtx, reflectpb.NewServerReflectionClient(cc))

	// descSource, _ := getDescSource(ctx, refClient)

	// log.Default().Println(descSource)
	// arrange for the RPCs to be cleanly shutdown
	// reset := func() {
	// 	if refClient != nil {
	// 		refClient.Reset()
	// 		refClient = nil
	// 	}
	// 	if cc != nil {
	// 		cc.Close()
	// 		cc = nil
	// 	}
	// }
	// defer reset()

	// exit = func(code int) {
	// 	// since defers aren't run by os.Exit...
	// 	reset()
	// 	os.Exit(code)
	// }

	return refClient
}

func Dial(ctx context.Context, config diaConfig) *grpc.ClientConn {
	target := config.target
	// 超时时间
	dialTime := 10 * time.Second
	if config.connectTimeout > 0 {
		dialTime = time.Duration(config.connectTimeout * float64(time.Second))
	}
	// 上下文
	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var opts []grpc.DialOption
	// 配置
	if config.keepaliveTime > 0 {
		timeout := time.Duration(config.keepaliveTime * float64(time.Second))
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    timeout,
			Timeout: timeout,
		}))
	}
	if config.maxMsgSz > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.maxMsgSz)))
	}
	var creds credentials.TransportCredentials
	opts = append(opts, grpc.WithUserAgent(config.ua))

	network := "tcp"
	if config.isUnixSocket {
		network = "unix"
	}
	cc, err := dialer.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		fail(err, "Failed to dial target host %q", target)
	}
	return cc
}

func fail(err error, msg string, args ...interface{}) {
	if err != nil {
		msg += ": %v"
		args = append(args, err)
	}
	fmt.Fprintf(os.Stderr, msg, args...)
	fmt.Fprintln(os.Stderr)
}
