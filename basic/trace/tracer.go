package trace

import (
	"fmt"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

//目前仅支持jeager

type TraceService int32

const (
	NONE TraceService = iota
	JAEGER
)

func InitOpentracing(serviceName, endpoint string) {
	cfg := jaegercfg.Configuration{
		// 将采样频率设置为 1，每一个 span 都记录，方便查看测试结果
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: false,
			// 将 span 发往 jaeger-collector 的服务地址
			CollectorEndpoint: fmt.Sprintf("http://%s/api/traces", endpoint),
		},
	}
	_, err := cfg.InitGlobalTracer(serviceName, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return
}
