package serverinterceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryCrashInterceptor(t *testing.T) {
	assert.NotPanics(t, func() {
		UnaryCrashInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/Crash",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("crash")
		})
	})
}

func TestStreamCrashInterceptor(t *testing.T) {
	assert.NotPanics(t, func() {
		StreamCrashInterceptor(nil, nil, &grpc.StreamServerInfo{
			FullMethod: "/StreamCrash",
		}, func(srv interface{}, stream grpc.ServerStream) error {
			panic("stream crash")
		})
	})
}
