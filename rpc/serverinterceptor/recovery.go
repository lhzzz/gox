package serverinterceptor

import (
	"context"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"singer.com/basic/log"
)

// StreamCrashInterceptor catches panics in processing stream requests and recovers.
func StreamCrashInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(stream.Context(), r)
	})

	return handler(svr, stream)
}

// UnaryCrashInterceptor catches panics in processing unary requests and recovers.
func UnaryCrashInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(ctx, r)
	})

	return handler(ctx, req)
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(ctx context.Context, r interface{}) error {
	const size = 64 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	if _, ok := r.(string); ok {
		log.RpcPanicf(ctx, "Observed a panic: %s\n%s", r, stacktrace)
	} else {
		log.RpcPanicf(ctx, "Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
	return status.Error(codes.Internal, "panic")
}
