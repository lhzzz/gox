package serverinterceptor

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"singer.com/basic/meta"
	"singer.com/util/color"
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
		requestId := meta.GetReuqestId(ctx)
		content, err := json.Marshal(req)
		if err != nil {
			logrus.WithField("RequestId", requestId).Errorf("%s - %s", addr, err.Error())
		} else if duration > time.Duration(slowThreshold) {
			logrus.WithField("Cost", duration).WithField("RequestId", requestId).Warnf("%s %s -> %s - %s", color.YollowStr("[RPC-SlowCall]"), addr, method, string(content))
		}
	}
}
