package serverinterceptor

import (
	"singer.com/basic/trace"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func UnaryOpentracingInterceptor() grpc.UnaryServerInterceptor {
	return otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.IncludingSpans(otgrpc.SpanInclusionFunc(trace.SpanInclusionFunc)),
		otgrpc.SpanDecorator(otgrpc.SpanDecoratorFunc(trace.SpanDecoratorError)))
}

// OpentracingGrpcStreamServerInterceptor return opentracing grpc stream server interceptor
func StreamOpentracingInterceptor() grpc.StreamServerInterceptor {
	return otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer(),
		otgrpc.IncludingSpans(otgrpc.SpanInclusionFunc(trace.SpanInclusionFunc)),
		otgrpc.SpanDecorator(otgrpc.SpanDecoratorFunc(trace.SpanDecoratorError)))
}
