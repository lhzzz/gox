package clientinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"singer.com/basic/breaker"
)

func UnaryBreakerInterceptor(bkr breaker.Breaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return bkr.Do(method, func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}

func StreamBreakerInterceptor(bkr breaker.Breaker) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
		err = bkr.Do(method, func() error {
			var err error
			cs, err = streamer(ctx, desc, cc, method, opts...)
			return err
		})
		return cs, err
	}
}
