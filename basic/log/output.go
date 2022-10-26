package log

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"singer.com/basic/meta"
	"singer.com/util/color"
)

const (
	costKey      = "Cost"
	requestIdKey = "RequestId"

	rpcSlowCallKey = "[RPC-SlowCall] "
	rpcErrorKey    = "[RPC-ERR] "
	rpcPanicKey    = "[RPC-PANIC] "
)

func RpcSlowf(ctx context.Context, d time.Duration, format string, args ...interface{}) {
	requestId := meta.GetReuqestId(ctx)
	logrus.WithField(costKey, d).WithField(requestIdKey, requestId).Warnf(color.Yellow(rpcSlowCallKey)+format, args...)
}

func RpcErrorf(ctx context.Context, format string, args ...interface{}) {
	requestId := meta.GetReuqestId(ctx)
	logrus.WithField(requestIdKey, requestId).Errorf(color.Red(rpcErrorKey)+format, args...)
}

//will not panic, only log
func RpcPanicf(ctx context.Context, format string, args ...interface{}) {
	requestId := meta.GetReuqestId(ctx)
	logrus.WithField(requestIdKey, requestId).Errorf(color.Red(rpcPanicKey)+format, args...)
}
