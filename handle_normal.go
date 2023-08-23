package gogrpcgeneric

import (
	"context"
	"io"

	//lint:ignore SA1019 we have to import this because it appears in exported API
	"github.com/golang/protobuf/proto" //lint:ignore SA1019 we have to import this because it appears in exported API
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// -------------------------------------------------------------------------
// Handle 用来处理chan 消息
// grpc 的返回值需要通过chan 返回给 application'user
// -------------------------------------------------------------------------

type Response struct {
	Result string
	err    error
}

type NormalEventHandler struct {
	Out       io.Writer
	Formatter Formatter
	// 0 = default
	// 1 = verbose
	// 2 = very verbose
	VerbosityLevel int

	// NumResponses is the number of responses that have been received.
	NumResponses int
	// Status is the status that was received at the end of an RPC. It is
	// nil if the RPC is still in progress.
	Status *status.Status
	// 核心响应
	Response chan Response
}

// NewDefaultEventHandler returns an InvocationEventHandler that logs events to
// the given output. If verbose is true, all events are logged. Otherwise, only
// response messages are logged.
//
// Deprecated: NewDefaultEventHandler exists for compatibility.
// It doesn't allow fine control over the `VerbosityLevel`
// and provides only 0 and 1 options (which corresponds to the `verbose` argument).
// Use DefaultEventHandler{} initializer directly.
func NewNormalEventHandler(out io.Writer, descSource DescriptorSource, formatter Formatter, verbose bool) *DefaultEventHandler {
	verbosityLevel := 0
	if verbose {
		verbosityLevel = 1
	}
	return &DefaultEventHandler{
		Out:            out,
		Formatter:      formatter,
		VerbosityLevel: verbosityLevel,
	}
}

var _ InvocationEventHandler = (*NormalEventHandler)(nil)

func (h *NormalEventHandler) OnResolveMethod(ctx context.Context, md *desc.MethodDescriptor) {
	if h.VerbosityLevel > 0 {
		_, err := GetDescriptorText(md, nil)
		if err == nil {
		}
	}
}

func (h *NormalEventHandler) OnSendHeaders(ctx context.Context, md metadata.MD) {
}

func (h *NormalEventHandler) OnReceiveHeaders(ctx context.Context, md metadata.MD) {
}

func (h *NormalEventHandler) OnReceiveResponse(ctx context.Context, resp proto.Message) {
	if h.Response == nil {
		h.Response = make(chan Response, 1)
	}
	h.NumResponses++
	if h.VerbosityLevel > 1 {
	}
	if h.VerbosityLevel > 0 {
	}
	if respStr, err := h.Formatter(resp); err != nil {
		h.Response <- Response{
			Result: respStr,
			err:    err,
		}
	} else {
		h.Response <- Response{
			Result: respStr,
			err:    err,
		}
	}
}

func (h *NormalEventHandler) OnReceiveTrailers(ctx context.Context, stat *status.Status, md metadata.MD) {
	h.Status = stat
}
