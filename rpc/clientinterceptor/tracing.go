package clientinterceptor

import (
	"context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	eventKey   = "event"
	messageKey = "message"
)

func TraceInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx, span := startSpan(ctx, method, cc.Target())
	defer span.Finish()

	trailer := metadata.MD{}
	opts = append(opts, grpc.Trailer(&trailer))
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String(eventKey, "error"), log.String(messageKey, err.Error()))
	}
	return err
}

func startSpan(ctx context.Context, method, target string) (context.Context, opentracing.Span) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, method)
	return ctx, span
}
