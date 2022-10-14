package micro

import (
	"math"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc/keepalive"
	"singer.com/basic/limit"
)

var defaultKaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var deafultKasp = keepalive.ServerParameters{
	MaxConnectionAgeGrace: 5 * time.Second, // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  2 * time.Minute, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               5 * time.Second, // Wait 5 second for the ping ack before assuming the connection is dead
}

type Options struct {
	serverName            string        //服务名
	listenOn              string        //服务启动地址
	metricListenAddr      string        //统计监听地址
	metricName            string        //统计名
	enableMetric          bool          //使能统计
	logListenAddr         string        //动态修改日志级别的监听地址
	enableLogServer       bool          //使能日志监听
	slowThreshold         time.Duration //慢日志阈值
	openTraceAddress      string        //调用链服务地址
	pprofListenAddr       string        //pprof监听地址
	enablePProf           bool          //使能pprof
	timeout               int64         //超时退出机制
	maxRecvMsgSize        int           //设置rpc所能接受的最大消息长度
	maxSendMsgSize        int           //设置rpc所能发送的最大消息长度
	enableKeepAlivePolicy bool          //使能keepalive EnforcementPolicy
	kaep                  keepalive.EnforcementPolicy
	kasp                  keepalive.ServerParameters
	limiter               limit.Limiter  //限流器
	preRunHooks           []func() error //服务启动前需要执行的操作
	preShutdownHooks      []func() error //服务停止时需要执行的操作
}

type Option func(*Options)

func defaultOption() Options {
	return Options{
		listenOn:              ":50051",
		metricListenAddr:      ":9090",
		logListenAddr:         ":1065",
		pprofListenAddr:       ":1066",
		openTraceAddress:      viper.GetString("JAEGER.ADDR"),
		enableLogServer:       true,
		maxRecvMsgSize:        4 * 1024 * 1024, //grpc default recv msg size
		maxSendMsgSize:        math.MaxInt32,   // grpc default send msg size
		enableKeepAlivePolicy: false,
		kaep:                  defaultKaep,
		kasp:                  deafultKasp,
	}
}

func newOptions(opts ...Option) Options {
	opt := defaultOption()
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Name(name string) Option {
	return func(o *Options) {
		o.serverName = name
	}
}

func Address(addr string) Option {
	return func(o *Options) {
		o.listenOn = addr
	}
}

func SlowThreshold(duration time.Duration) Option {
	return func(o *Options) {
		o.slowThreshold = duration
	}
}

func EnableMetric(metricName string) Option {
	return func(o *Options) {
		o.enableMetric = true
		o.metricName = metricName
	}
}

func EnablePProf(addr string) Option {
	return func(o *Options) {
		o.enablePProf = true
		o.pprofListenAddr = addr
	}
}

func SetMaxRecvMsgSize(size int) Option {
	return func(o *Options) {
		o.maxRecvMsgSize = size
	}
}

func SetMaxSendMsgSize(size int) Option {
	return func(o *Options) {
		o.maxSendMsgSize = size
	}
}

func SetMaxConnectionIdle(d time.Duration) Option {
	return func(o *Options) {
		o.kasp.MaxConnectionIdle = d
	}
}

func SetKeepAliveEnforcementPolicy(kaep keepalive.EnforcementPolicy) Option {
	return func(o *Options) {
		o.enableKeepAlivePolicy = true
		o.kaep = kaep
	}
}

func Timeout(timeout int64) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func Limiter(l limit.Limiter) Option {
	return func(o *Options) {
		o.limiter = l
	}
}
