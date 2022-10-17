package clientinterceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestTimeoutInterceptor(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			"a1",
			time.Second * 5,
		},
		{
			"a2",
			time.Second * 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cc := new(grpc.ClientConn)
			interceptor := TimeoutInterceptor(3 * time.Second)
			err := interceptor(context.Background(), "/foo", nil, nil, cc, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
				opts ...grpc.CallOption) error {
				timer := time.NewTimer(test.duration)
				defer timer.Stop()
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-timer.C:
					t.Log("success")
				}
				return nil
			})
			if test.duration > 3*time.Second {
				assert.EqualError(t, err, context.DeadlineExceeded.Error())
			} else {
				assert.Nil(t, err)
			}

		})
	}
}
