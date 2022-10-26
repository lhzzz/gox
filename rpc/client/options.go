package client

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"singer.com/basic/breaker"
	"singer.com/rpc/clientinterceptor"
)

type ClientOptions struct {
	block         bool                             //连接是否阻塞
	dialOptions   []grpc.DialOption                //grpc 连接options
	enableTrace   bool                             //分布式调用链追踪
	enableMeta    bool                             //元数据携带
	breaker       breaker.Breaker                  //熔断器
	timeout       time.Duration                    //超时调用
	slowThreshold time.Duration                    //慢日志阈值
	maxMsgSize    int                              //最大消息大小
	creds         credentials.TransportCredentials //连接证书
	retryConf     *clientinterceptor.RetryConfigs  //重试配置
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

func WithBreaker(bkr breaker.Breaker) ClientOption {
	return func(co *ClientOptions) {
		co.breaker = bkr
	}
}

func WithSlowThreshold(d time.Duration) ClientOption {
	return func(co *ClientOptions) {
		co.slowThreshold = d
	}
}

func WithTimeout(d time.Duration) ClientOption {
	return func(co *ClientOptions) {
		co.timeout = d
	}
}

func WithBlock() ClientOption {
	return func(co *ClientOptions) {
		co.block = true
	}
}

func WithMaxMsgSize(max int) ClientOption {
	return func(co *ClientOptions) {
		co.maxMsgSize = max
	}
}

func WithCreds(creds credentials.TransportCredentials) ClientOption {
	return func(co *ClientOptions) {
		co.creds = creds
	}
}

func WithRetry(retries map[string]clientinterceptor.RetryConfig) ClientOption {
	return func(co *ClientOptions) {
		co.retryConf = &clientinterceptor.RetryConfigs{Configs: retries}
	}
}
