package clientinterceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestSlowlogInterceptor(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "hello",
			err:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cc := new(grpc.ClientConn)
			err := UnarySlowlogInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					time.Sleep(time.Second * 3)
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}
}
