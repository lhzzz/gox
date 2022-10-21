package serverinterceptor

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"singer.com/basic/errorx"
)

func TestUnaryErrorInterceptor(t *testing.T) {
	resp, err := UnaryErrorInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/Unary/Error",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errorx.Wrapf(errorx.TOKEN_GENERATE_ERROR, "token gen failed :%s", "db error")
	})

	t.Log(resp, err)
}

func TestStreamErrorInterceptor(t *testing.T) {
	err := StreamErrorInterceptor(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
		FullMethod: "/Stream/Error",
	}, func(srv interface{}, stream grpc.ServerStream) error {
		return errorx.Wrap(errorx.SERVER_COMMON_ERROR, "server internal error")
	})

	t.Log(err)
}
