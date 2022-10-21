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
