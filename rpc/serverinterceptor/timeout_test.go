package serverinterceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryTimeoutInterceptor(t *testing.T) {
	intercept := UnaryTimeoutInterceptor(3 * time.Second)

	_, err := intercept(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/Timeout",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		for i := 0; i < 10; i++ {
			t.Log("not graceful", i)
			time.Sleep(1 * time.Second) //this is not graceful, the handler will go on even the intercept has return.
		}
		t.Log("not graceful finish")
		return nil, nil
	})
	assert.Equal(t, err, status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error()))

	_, err = intercept(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/Timeout",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		done := false
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				done = true
				break //check the context whether if done
			default:
				t.Log("graceful", i)
				time.Sleep(1 * time.Second)
			}
		}
		assert.EqualValues(t, done, true)
		t.Log("graceful finish")
		return nil, nil
	})

	time.Sleep(1 * time.Second)
}
