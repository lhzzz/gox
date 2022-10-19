package clientinterceptor

import (
	"context"
	"path"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const defaultSlowThreshold int64 = int64(time.Millisecond * 500)

var (
	slowThreshold = defaultSlowThreshold
)

func SetSlowThreshold(d time.Duration) {
	atomic.StoreInt64(&slowThreshold, d.Nanoseconds())
}

// SlowlogInterceptor is an interceptor that logs the processing time.
func UnarySlowlogInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Since(start)
	if elapsed > time.Duration(slowThreshold) {
		logrus.WithContext(ctx).Warnf("[RPC] [Cost:%v] ok - slowcall - %s - %v - %v",
			elapsed, serverName, req, reply)
	}
	return err
}
