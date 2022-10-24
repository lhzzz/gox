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
	var errBadRequest = fmt.Errorf("bad request")
	var errInternal = fmt.Errorf("internal error")

	b := breaker.NewPrefixHystrixBreaker("/api", breaker.NewHystrixConfig(4, 50, func(err error) bool {
		if err == errBadRequest { //bad request不计入熔断的统计
			return true
		}
		return false
	}, func(err error) error {
		t.Log("fallback", err)
		return nil
	}))

	interceptor := UnaryBreakerInterceptor(b)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api1",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errBadRequest
		})
		t.Log(resp, err)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api2",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errBadRequest
		})
		t.Log(resp, err)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api3",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errInternal
		})
		t.Log(resp, err)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api4",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errInternal
		})
		t.Log(resp, err)
	}

	//达到50%了，熔断
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/api5",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errInternal
		})
		t.Log(resp, err)
	}
}

func TestStreamBreakerInterceptor(t *testing.T) {
	var errBadRequest = fmt.Errorf("bad request")
	var errInternal = fmt.Errorf("internal error")

	b := breaker.NewMutiHystrixBreaker(map[string]breaker.HystrixConfig{
		"/api/CreateStream": breaker.NewHystrixConfig(4, 50, func(err error) bool {
			if err == errBadRequest {
				return true
			}
			return false
		}, func(err error) error {
			t.Log("fallback handle err:", err)
			return nil
		}),
	})

	it := StreamBreakerInterceptor(b)

	//accept, not metric error
	for i := 0; i < 4; i++ {
		err := it(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
			FullMethod: "/api/CreateStream",
		}, func(srv interface{}, stream grpc.ServerStream) error {
			return errBadRequest
		})
		t.Log(err)
	}

	//not hit method, not metric error
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		err := it(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
			FullMethod: "/api/UpdateStream",
		}, func(srv interface{}, stream grpc.ServerStream) error {
			return errInternal
		})
		t.Log(err)
	}

	//metric error
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		err := it(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
			FullMethod: "/api/CreateStream",
		}, func(srv interface{}, stream grpc.ServerStream) error {
			return errInternal
		})
		t.Log(err)
	}

	//circuit open
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		err := it(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
			FullMethod: "/api/CreateStream",
		}, func(srv interface{}, stream grpc.ServerStream) error {
			t.Log("I will success")
			return nil
		})
		t.Log(err)
	}
}
