package trace

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

const (
	healthCheck = "/grpc.health.v1.Health/Check"
)

func checkMethod(method string) bool {
	return method != healthCheck
}

// 目前只针对方法做了过滤，多个过滤条件使用交集
func SpanInclusionFunc(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
	return checkMethod(method)
}

func SpanDecoratorError(span opentracing.Span, method string, req, resp interface{}, grpcError error) {
	if grpcError != nil {
		span.SetTag("error", true)
		span.LogFields(
			opentracing_log.String("event", "error"),
			opentracing_log.String("message", fmt.Sprintf("%+v", grpcError)),
		)
	}
}
