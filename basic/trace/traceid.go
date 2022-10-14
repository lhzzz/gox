package trace

import (
	"context"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

type TraceID interface {
	// judge traceID whether it is valid
	IsValid() bool
	// traceID output format
	String() string
}

func traceIDConversionRequestID(traceID TraceID) string {
	if !traceID.IsValid() {
		return uuid.New().String()
	}
	return traceID.String()
}

func TraceIDFromContext(ctx context.Context) TraceID {
	return jaegerTraceIDFromContext(ctx)
}

func jaegerTraceIDFromContext(ctx context.Context) TraceID {
	sp := opentracing.SpanFromContext(ctx)
	if sp == nil {
		return jaeger.TraceID{}
	}
	return sp.Context().(jaeger.SpanContext).TraceID()
}
