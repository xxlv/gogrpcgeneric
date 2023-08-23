package gogrpcgeneric

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"google.golang.org/grpc"
)

// InvokerPayload Generic invoke Payload
type InvokerPayload struct {
	Service   string
	Method    string
	JsonParam string
	Metadata  []string // [header_key:value1] [header_key:value2] are both supported
	Handler   InvocationEventHandler
}

// GenericClient Generic client
type GenericClient struct {
	ConnectTimeout float64
	KeepaliveTime  float64
	Registryconfig *NacosRegistryConfig
	// generic invoke group
	Group  string
	Logger Logger
	Debug  bool
}

func NewGenericClient() *GenericClient {
	c := &GenericClient{}
	if c.Logger == nil {
		c.Logger = &DefaultLogger{}
	}
	return c
}

// GenericUnaryInvoke.
func (gc *GenericClient) GenericUnaryInvoke(ctx context.Context, registryService string, payload InvokerPayload) chan Response {
	hd := &NormalEventHandler{}
	payload.Handler = hd
	gc.genericInvoke(ctx, registryService, payload.Service, payload.Method, payload.Metadata, payload.JsonParam, payload.Handler)
	return hd.Response
}

// GenericUnaryInvokeAsync.
// Note: now only support unary invoke
func (gc *GenericClient) GenericUnaryInvokeAsync(ctx context.Context, registryService string, payload InvokerPayload) {
	if payload.Handler == nil {
		gc.Logger.Info("generic invoke async current request miss handle , will ignore response")
	}
	gc.genericInvoke(ctx, registryService, payload.Service, payload.Method, payload.Metadata, payload.JsonParam, payload.Handler)
}

func (gc *GenericClient) GetReflectClient(ctx context.Context, service string) (*grpc.ClientConn, *grpcreflect.Client, error) {
	if gc.ConnectTimeout <= 0 {
		gc.ConnectTimeout = 1000 //ms
	}
	if gc.KeepaliveTime <= 0 {
		gc.KeepaliveTime = 1000 //ms
	}
	var ins *model.Instance
	group := gc.Group
	if gc.Registryconfig == nil {
		ins, _ = LoadServiceDefault(service, group)
	} else {
		ins, _ = LoadServiceLocation(*gc.Registryconfig, service, group)
	}

	if ins == nil {
		return nil, nil, errors.New("can not found the service instance")
	}

	cc := Dial(ctx, diaConfig{
		connectTimeout: gc.ConnectTimeout,
		keepaliveTime:  gc.KeepaliveTime,
		target:         fmt.Sprintf("%s:%s", ins.Ip, fmt.Sprint(ins.Port)),
	})
	refClient := doGetReflectClient(ctx, cc)

	// Note: if gc happend, stream will auto closed
	return cc, refClient, nil
}

// Invoke.
func (gc *GenericClient) genericInvoke(ctx context.Context, registryService string, service string, method string, headers []string, jsondata string, handle InvocationEventHandler) error {
	serviceSymbol := fmt.Sprintf("%s.%s", service, method)
	if gc.Debug {
		gc.Logger.Debugf("generic invoke with service symbol %s", serviceSymbol)
	}
	cc, refClient, err := gc.GetReflectClient(ctx, registryService)
	if err != nil {
		return err
	}
	ds, _ := GetDescSource(ctx, refClient)
	in := strings.NewReader(jsondata)

	options := FormatOptions{
		EmitJSONDefaultFields: true,
		IncludeTextSeparator:  true,
		AllowUnknownFields:    true,
	}
	rf, formatter, _ := RequestParserAndFormatter(Format("json"), ds, in, options)

	var h InvocationEventHandler = handle
	verbosityLevel := 1
	if h == nil {
		gc.Logger.Info("event handle is miss , choice default handle ,verbosityLevel is ", verbosityLevel)

		h = &DefaultEventHandler{
			Out:            os.Stdout,
			VerbosityLevel: verbosityLevel,
			Formatter:      formatter,
		}
	}

	// NormalEventHandle没有formatter
	if v, ok := h.(*NormalEventHandler); ok {
		if v.Formatter == nil {
			v.Formatter = formatter
		}
		if v.Out == nil {
			v.Out = os.Stdout
		}
		if v.VerbosityLevel <= 0 {
			v.VerbosityLevel = verbosityLevel
		}
	}
	// do invoke PRC
	return InvokeRPC(ctx, ds, cc, serviceSymbol, headers, h, rf.Next)
}
