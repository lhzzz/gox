package clientinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"singer.com/basic/breaker"
)

func UnaryBreakerInterceptor(bkr breaker.Breaker, accept breaker.Acceptable) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return bkr.DoWithAcceptable(func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}, accept)
	}
}
