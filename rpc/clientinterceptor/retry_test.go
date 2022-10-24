package clientinterceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"singer.com/basic/errorx"
)

func TestWaitRetryBackoff(t *testing.T) {
	rc := &RetryConfig{
		max:            3,
		perCallTimeout: 3 * time.Second,
		includeHeader:  true,
		codes:          []errorx.ErrorCode{errorx.TOKEN_EXPIRE_ERROR},
		backoffFunc: func(ctx context.Context, attempt uint) time.Duration {
			return time.Second * 3
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("requestId", "ef490lpdav4"))
	waitRetryBackoff(1, ctx, rc)

}

func TestContextErrToGrpcErr(t *testing.T) {
	assert.Equal(t, status.Code(contextErrToGrpcErr(context.Canceled)), codes.Canceled)
	assert.Equal(t, status.Code(contextErrToGrpcErr(context.DeadlineExceeded)), codes.DeadlineExceeded)
	assert.Equal(t, status.Code(errorx.Wrap(errorx.DB_ERROR, "db conn failed")), codes.Unknown)
}

func TestPerCallContext(t *testing.T) {
	rc := &RetryConfig{
		max:            3,
		perCallTimeout: 3 * time.Second,
		includeHeader:  true,
		codes:          []errorx.ErrorCode{errorx.TOKEN_EXPIRE_ERROR},
		backoffFunc: func(ctx context.Context, attempt uint) time.Duration {
			return time.Second * 3
		},
	}

	pctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < int(rc.max); i++ {
		ctx := perCallContext(pctx, rc, 1)
		t.Log(ctx)
	}
}

func TestIsContextError(t *testing.T) {
	assert.True(t, isContextError(status.Error(codes.Canceled, "cancel")))
	assert.True(t, isContextError(status.Error(codes.DeadlineExceeded, "deadline")))
	assert.False(t, isContextError(status.Error(codes.Aborted, "abort")))
}

func TestIsRetriable(t *testing.T) {
	rc := &RetryConfig{
		max:            3,
		perCallTimeout: 3 * time.Second,
		includeHeader:  true,
		codes:          []errorx.ErrorCode{errorx.TOKEN_EXPIRE_ERROR, errorx.ErrorCode(codes.AlreadyExists)},
		backoffFunc: func(ctx context.Context, attempt uint) time.Duration {
			return time.Second * 3
		},
	}

	assert.True(t, isRetriable(status.Error(codes.Code(errorx.TOKEN_EXPIRE_ERROR), "token expire"), rc))
	assert.True(t, isRetriable(status.Error(codes.AlreadyExists, "exist"), rc))
	assert.False(t, isRetriable(status.Error(codes.Canceled, "cancel"), rc))
	assert.False(t, isRetriable(status.Error(codes.Code(errorx.DB_ERROR), "db error"), rc))
	assert.False(t, isRetriable(status.Error(codes.DeadlineExceeded, "deadline"), rc))
}

func TestUnaryRetryInterceptor(t *testing.T) {
	confgis := RetryConfigs{
		Configs: map[string]RetryConfig{
			"/api/Create": {
				max:            3,
				perCallTimeout: 3 * time.Second,
				includeHeader:  true,
				codes:          []errorx.ErrorCode{errorx.TOKEN_EXPIRE_ERROR},
				backoffFunc: func(ctx context.Context, attempt uint) time.Duration {
					return time.Second * 3
				},
			},
		},
	}

	interceptor := UnaryRetryInterceptor(confgis)

	pctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	cc := new(grpc.ClientConn)
	err := interceptor(pctx, "/api/Create", nil, nil, cc, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		opts ...grpc.CallOption) error {
		time.Sleep(time.Second * 3)
		return status.Error(codes.Code(errorx.TOKEN_EXPIRE_ERROR), "")
	})
	t.Log(err)
}
