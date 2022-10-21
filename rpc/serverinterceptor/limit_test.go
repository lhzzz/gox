package serverinterceptor

import (
	"context"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"singer.com/basic/limit"
)

func TestUnaryLimitInterceptor(t *testing.T) {
	l := limit.NewPrefixTokenLimiter("/api", limit.Capcaity(5), limit.FillInterval(time.Second), limit.Quantum(1))

	interceptor := UnaryLimitInterceptor(l)

	for i := 0; i < 10; i++ {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
				FullMethod: "/api/Unary",
			}, func(ctx context.Context, req interface{}) (interface{}, error) {
				time.Sleep(1 * time.Second)
				return nil, nil
			})
			t.Log(resp, err)
		})
	}
}

func TestStreamLimitInterceptor(t *testing.T) {
	l := limit.NewMutiTokenLimiter(map[string][]limit.LimitOption{
		"/api/Stream": {limit.Capcaity(5), limit.FillInterval(1 * time.Second), limit.Quantum(1)},
	})

	var wg sync.WaitGroup
	interceptor := StreamLimitInterceptor(l)
	for i := 0; i < 20; i++ {
		if i < 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := interceptor(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
					FullMethod: "/api/Stream",
				}, func(srv interface{}, stream grpc.ServerStream) error {
					time.Sleep(1 * time.Second)
					return nil
				})
				t.Log(err)
			}()
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := interceptor(nil, mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
					FullMethod: "/api/Stream2",
				}, func(srv interface{}, stream grpc.ServerStream) error {
					time.Sleep(1 * time.Second)
					return nil
				})
				t.Log(err)
			}()
		}
	}
	wg.Wait()
}
