package serverinterceptor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"singer.com/basic/breaker"
)

func TestUnaryBreakerInterceptor(t *testing.T) {
	b := breaker.NewPrefixHystrixBreaker("/api", breaker.NewHystrixConfig(4, 50, nil, func(err error) error {
		t.Log("fallback", err)
		return nil
	}))

	interceptor := UnaryBreakerInterceptor(b)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, fmt.Errorf("bad request")
		})
		t.Log(resp, err)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, fmt.Errorf("bad request")
		})
		t.Log(resp, err)
	}

}

func TestStreamBreakerInterceptor(t *testing.T) {

}
