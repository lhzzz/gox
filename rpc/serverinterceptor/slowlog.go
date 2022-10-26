package serverinterceptor

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"singer.com/basic/log"
)

const defaultSlowThreshold int64 = int64(time.Millisecond * 500)

var (
	notLoggingContentMethods sync.Map
	slowThreshold            = defaultSlowThreshold
)

func SetSlowThreshold(d time.Duration) {
	atomic.StoreInt64(&slowThreshold, d.Nanoseconds())
}

// UnaryStatInterceptor returns a func that uses given metrics to report stats.
func UnarySlowlogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			logDuration(ctx, info.FullMethod, req, duration)
		}()
		return handler(ctx, req)
	}
}

func logDuration(ctx context.Context, method string, req interface{}, duration time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		addr = client.Addr.String()
	}

	_, ok = notLoggingContentMethods.Load(method)
	if !ok {
		content, err := json.Marshal(req)
		if err != nil {
			log.RpcErrorf(ctx, "%s - %s", addr, err.Error())
		} else if duration > time.Duration(slowThreshold) {
			log.RpcSlowf(ctx, duration, "%s -> %s - %s", addr, method, string(content))
		}
	}
}
