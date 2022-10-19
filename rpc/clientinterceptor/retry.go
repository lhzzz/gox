package clientinterceptor

import (
	"context"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"singer.com/basic/errorx"
)

const (
	kAttemptMetadataKey = "x-retry-attempty"
)

type BackoffFunc func(attempt uint) time.Duration
type BackoffFuncContext func(ctx context.Context, attempt uint) time.Duration

type RetryConfig struct {
	max            uint
	perCallTimeout time.Duration
	includeHeader  bool
	codes          []errorx.ErrorCode
	backoffFunc    BackoffFuncContext
}

type RetryConfigs struct {
	Configs map[string]RetryConfig
}

func UnaryRetryInterceptor(rc RetryConfigs) grpc.UnaryClientInterceptor {
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if rc.Configs == nil {
			return invoker(parentCtx, method, req, reply, cc, opts...)
		}
		if config, ok := rc.Configs[method]; !ok {
			return invoker(parentCtx, method, req, reply, cc, opts...)
		} else {
			if config.max == 0 {
				return invoker(parentCtx, method, req, reply, cc, opts...)
			}
			var lastErr error
			for attempt := uint(0); attempt < config.max; attempt++ {
				if err := waitRetryBackoff(attempt, parentCtx, &config); err != nil {
					return err
				}
				callCtx := perCallContext(parentCtx, &config, attempt)
				lastErr = invoker(callCtx, method, req, reply, cc, opts...)
				// TODO: Maybe dial and transport errors should be retriable?
				if lastErr == nil {
					return nil
				}
				logrus.Error(parentCtx, "grpc_retry attempt: %d, got err: %v", attempt, lastErr)
				if isContextError(lastErr) {
					if parentCtx.Err() != nil {
						logrus.Error(parentCtx, "grpc_retry attempt: %d, parent context error: %v", attempt, parentCtx.Err())
						// its the parent context deadline or cancellation.
						return lastErr
					} else if config.perCallTimeout != 0 {
						// We have set a perCallTimeout in the retry middleware, which would result in a context error if
						// the deadline was exceeded, in which case try again.
						logrus.Error(parentCtx, "grpc_retry attempt: %d, context error from retry call", attempt)
						continue
					}
				}
				if !isRetriable(lastErr, &config) {
					return lastErr
				}
			}
		}
		return nil
	}
}

func waitRetryBackoff(attempt uint, parentCtx context.Context, rc *RetryConfig) error {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = rc.backoffFunc(parentCtx, attempt)
	}
	if waitTime > 0 {
		logrus.Info(parentCtx, "grpc_retry attempt: %d, backoff for %v", attempt, waitTime)
		timer := time.NewTimer(waitTime)
		select {
		case <-parentCtx.Done():
			timer.Stop()
			return contextErrToGrpcErr(parentCtx.Err())
		case <-timer.C:
		}
	}
	return nil
}

func contextErrToGrpcErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

func perCallContext(parentCtx context.Context, rc *RetryConfig, attempt uint) context.Context {
	ctx := parentCtx
	if rc.perCallTimeout != 0 {
		ctx, _ = context.WithTimeout(ctx, rc.perCallTimeout)
	}
	if attempt > 0 && rc.includeHeader {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		md.Set(kAttemptMetadataKey, strconv.FormatUint(uint64(attempt), 10))
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func isContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}

func isRetriable(err error, rc *RetryConfig) bool {
	if isContextError(err) {
		// context errors are not retriable based on user settings.
		return false
	}
	errCode := errorx.ErrorCode(status.Code(err))
	for _, code := range rc.codes {
		if code == errCode {
			return true
		}
	}
	return false
}
