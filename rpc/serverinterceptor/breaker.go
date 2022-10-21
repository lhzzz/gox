package serverinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"singer.com/basic/breaker"
)

// UnaryLimitInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryBreakerInterceptor(bkr breaker.Breaker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = bkr.Do(info.FullMethod, func() error {
			var herr error
			resp, herr = handler(ctx, req)
			return herr
		})
		return resp, err
	}
}

// StreamLimitInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamBreakerInterceptor(bkr breaker.Breaker) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		return bkr.Do(info.FullMethod, func() error {
			return handler(srv, stream)
		})
	}
}
