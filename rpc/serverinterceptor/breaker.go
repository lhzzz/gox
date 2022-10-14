package serverinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"singer.com/basic/breaker"
)

// UnaryLimitInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryBreakerInterceptor(bkr breaker.Breaker, accept breaker.Acceptable) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = bkr.DoWithAcceptable(func() error {
			var err error
			resp, err = handler(ctx, req)
			return err
		}, accept)
		return resp, err
	}
}

// StreamLimitInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamBreakerInterceptor(bkr breaker.Breaker, accept breaker.Acceptable) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = bkr.DoWithAcceptable(func() error {
			return handler(srv, stream)
		}, accept)
		return err
	}
}
