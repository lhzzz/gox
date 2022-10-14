package serverinterceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold)
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second.Nanoseconds(), slowThreshold)
}

func TestUnarySlowlogInterceptor(t *testing.T) {
	interceptor := UnarySlowlogInterceptor()
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/Test",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestLogDuration(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		req      interface{}
		duration time.Duration
	}{
		{
			name: "normal",
			ctx:  context.Background(),
			req:  "foo",
		},
		{
			name: "bad req",
			ctx:  context.Background(),
			req:  make(chan struct{}), // not marshalable
		},
		{
			name:     "timeout",
			ctx:      context.Background(),
			req:      "foo",
			duration: time.Second,
		},
		{
			name: "timeout",
			ctx:  context.Background(),
			req:  "foo",
		},
		{
			name:     "timeout",
			ctx:      context.Background(),
			req:      "foo",
			duration: time.Duration(slowThreshold) + time.Second,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.NotPanics(t, func() {
				logDuration(test.ctx, "foo", test.req, test.duration)
			})
		})
	}
}
