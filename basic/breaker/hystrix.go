package breaker

import (
	"strings"

	"singer.com/basic/hystrix"
)

const (
	// defaultTimeout is how long to wait for command to complete, in milliseconds
	defaultTimeout = 100000000
	// defaultMaxConcurrent is how many commands of the same type can run at the same time
	defaultMaxConcurrent = 100000000
	// defaultVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	defaultVolumeThreshold = 20
	// defaultSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	defaultSleepWindow = 5000
	// defaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	defaultErrorPercentThreshold = 50
)

type HystrixConfig struct {
	Cmd      hystrix.CommandConfig
	Accept   Acceptable        //判断错误是否需要计入错误统计中
	Fallback func(error) error // 失败处理逻辑，访问其他资源失败时，或者处于熔断开启状态时，会调用这段逻辑
}

type hystrixPrefixBreaker struct {
	Prefix string
	Config HystrixConfig
}

type hystrixMutiBreaker struct {
	Configs map[string]HystrixConfig
}

var (
	defaultHystrixCmdConfig = hystrix.CommandConfig{
		Timeout:                defaultTimeout,               // 接口请求的超时时间，单位ms，超过超时时间直接返回错误
		MaxConcurrentRequests:  defaultMaxConcurrent,         // 最大并发请求,超过并发直接返回错误
		SleepWindow:            defaultSleepWindow,           // 在熔断器被打开后，根据SleepWindow设置的时间控制多久后尝试服务是否可用，单位是ms
		RequestVolumeThreshold: defaultVolumeThreshold,       // 请求数量达到阈值后，会根据统计结果判断熔断
		ErrorPercentThreshold:  defaultErrorPercentThreshold, // 统计错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断，范围0-100
	}
)

// prefix： 服务名前缀   如果熔断某个服务某个类型接口，则填 /api.BackEndService/Create
//						如果熔断整个服务，则填 "/api.BackEndService"
//						如果熔断某块后台微服务，则填 "/api" (网关)
//						如果熔断整个后台，则填 "/"  (网关)
// conf： 熔断策略
func NewPrefixHystrixBreaker(prefix string, conf HystrixConfig) Breaker {
	hystrix.ConfigureCommand(prefix, conf.Cmd)
	return &hystrixPrefixBreaker{
		Prefix: prefix,
		Config: conf,
	}
}

// 	requestThreshold: 		请求数量达到阈值后，会根据统计结果判断熔断
// 	errorPercentThreshold:	统计错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
// 	accept： 				错误是否要被熔断统计,可为nil
// 	fallback: 				保底措施，触发熔断后执行这段逻辑，为nil则不执行,直接返回熔断错误(hystrix: circuit open)
func NewHystrixConfig(requestThreshold, errorPercentThreshold int, accept func(err error) bool, fallback func(error) error) HystrixConfig {
	cmd := defaultHystrixCmdConfig
	cmd.RequestVolumeThreshold = requestThreshold
	cmd.ErrorPercentThreshold = errorPercentThreshold
	return HystrixConfig{
		Accept:   accept,
		Fallback: fallback,
		Cmd:      cmd,
	}
}

func (hpb *hystrixPrefixBreaker) Do(method string, req func() error) error {
	if strings.HasPrefix(method, hpb.Prefix) {
		var retErr error
		var callback func() error
		var accept bool
		if hpb.Config.Accept != nil {
			callback = func() error {
				retErr = req()
				//如果设置了accept，且认为错误并不统计熔断中，则直接返回nil
				if retErr != nil {
					if hpb.Config.Accept(retErr) {
						accept = true
						return nil
					} else {
						return retErr
					}
				}
				return retErr
			}
		} else {
			callback = req
		}
		hyerr := hystrix.Do(hpb.Prefix, callback, hpb.Config.Fallback)
		if accept {
			hyerr = retErr
		}
		return hyerr
	} else {
		return req()
	}
}

//muti breaker
//is designed for muti rpc interface, like:
//	map[string]HystrixConfig{
//		"/api.BackEndService/Create": HystrixConfig1,
//		"/api.BackEndService/Update": HystrixConfig2,
//	}
func NewMutiHystrixBreaker(configs map[string]HystrixConfig) Breaker {
	for method, config := range configs {
		hystrix.ConfigureCommand(method, config.Cmd)
	}
	return &hystrixMutiBreaker{
		Configs: configs,
	}
}

func (hmb *hystrixMutiBreaker) Do(method string, req func() error) error {
	if config, ok := hmb.Configs[method]; ok {
		var retErr error
		var callback func() error
		var accept bool
		if config.Accept != nil {
			callback = func() error {
				retErr = req()
				if retErr != nil {
					if config.Accept(retErr) {
						accept = true
						return nil
					} else {
						return retErr
					}
				}
				return retErr
			}
		} else {
			callback = req
		}
		hyerr := hystrix.Do(method, callback, config.Fallback)
		if accept {
			hyerr = retErr
		}
		return hyerr
	} else {
		return req()
	}
}
