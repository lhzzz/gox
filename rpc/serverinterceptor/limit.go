package serverinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"singer.com/basic/limit"
)

// UnaryLimitInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryLimitInterceptor(limiter limit.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.LimitWithContext(ctx, info.FullMethod) {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// StreamLimitInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamLimitInterceptor(limiter limit.Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.LimitWithContext(stream.Context(), info.FullMethod) {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}
