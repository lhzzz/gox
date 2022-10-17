package client

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"singer.com/basic/breaker"
)

type ClientOptions struct {
	block          bool                             //连接是否阻塞
	dialOptions    []grpc.DialOption                //grpc 连接options
	enableTrace    bool                             //分布式调用链追踪
	enableMeta     bool                             //元数据携带
	breaker        breaker.Breaker                  //熔断器
	accept         breaker.Acceptable               //熔断依据
	timeout        time.Duration                    //超时调用
	slowThreshold  time.Duration                    //慢日志阈值
	maxRecvMsgSize int                              //最大接受消息大小
	maxSendMsgSize int                              //最大发送消息大小
	creds          credentials.TransportCredentials //连接证书
}

type ClientOption func(co *ClientOptions)

func newDefaultClientOptions() ClientOptions {
	return ClientOptions{
		enableTrace: true,
		creds:       insecure.NewCredentials(),
	}
}

func DialOptions(opts ...grpc.DialOption) ClientOption {
	return func(co *ClientOptions) {
		co.dialOptions = opts
	}
}

func DisableTrace() ClientOption {
	return func(co *ClientOptions) {
		co.enableTrace = false
	}
}

func EnableMeta() ClientOption {
	return func(co *ClientOptions) {
		co.enableMeta = true
	}
}

func WithBreaker(bkr breaker.Breaker, accept breaker.Acceptable) ClientOption {
	return func(co *ClientOptions) {
		co.breaker = bkr
		co.accept = accept
	}
}

func WithSlowThreshold(d time.Duration) ClientOption {
	return func(co *ClientOptions) {
		co.slowThreshold = d
	}
}

func WithTime(d time.Duration) ClientOption {
	return func(co *ClientOptions) {
		co.timeout = d
	}
}

func WithBlock() ClientOption {
	return func(co *ClientOptions) {
		co.block = true
	}
}

func WithMaxSendMsgSize(maxSend int) ClientOption {
	return func(co *ClientOptions) {
		co.maxSendMsgSize = maxSend
	}
}

func WithMaxRecvMsgSize(maxRecv int) ClientOption {
	return func(co *ClientOptions) {
		co.maxRecvMsgSize = maxRecv
	}
}

func WithCreds(creds credentials.TransportCredentials) ClientOption {
	return func(co *ClientOptions) {
		co.creds = creds
	}
}
