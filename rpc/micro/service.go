package micro

import (
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"singer.com/basic/log"
	"singer.com/basic/metric"
	"singer.com/basic/pprof"
	"singer.com/basic/trace"
	"singer.com/rpc/serverinterceptor"
	signalutil "singer.com/util/signal"
)

type service struct {
	app        Application
	opts       Options
	grpcServer *grpc.Server
	health     *health.Server
	stopCh     <-chan struct{}
}

func newService(app Application, opts ...Option) Service {
	service := new(service)
	options := newOptions(opts...)
	health := health.NewServer()
	health.SetServingStatus(options.serverName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	if options.enableLogServer {
		log.InitLogServer(options.logListenAddr)
	}

	if options.slowThreshold > 0 {
		serverinterceptor.SetSlowThreshold(options.slowThreshold)
	}

	if options.enableMetric {
		metric.InitPrometheusMetrics(options.metricName, options.metricListenAddr)
	}

	if options.enablePProf {
		pprof.Serve(options.pprofListenAddr)
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptor.UnaryGenerateMetadataInterceptor,
		serverinterceptor.UnarySlowlogInterceptor(),
		serverinterceptor.UnaryCrashInterceptor,
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		serverinterceptor.StreamGenerateMetadataInterceptor,
		serverinterceptor.StreamCrashInterceptor,
	}

	if options.timeout > 0 {
		unaryInterceptors = append(unaryInterceptors, serverinterceptor.UnaryTimeoutInterceptor(time.Duration(options.timeout)*time.Millisecond))
	}

	if options.limiter != nil {
		unaryInterceptors = append(unaryInterceptors, serverinterceptor.UnaryLimitInterceptor(options.limiter))
		streamInterceptors = append(streamInterceptors, serverinterceptor.StreamLimitInterceptor(options.limiter))
	}

	if options.breaker != nil && options.accecptable != nil {
		unaryInterceptors = append(unaryInterceptors, serverinterceptor.UnaryBreakerInterceptor(options.breaker, options.accecptable))
		streamInterceptors = append(streamInterceptors, serverinterceptor.StreamBreakerInterceptor(options.breaker, options.accecptable))
	}

	if len(options.openTraceAddress) > 0 {
		trace.InitOpentracing(options.serverName, options.openTraceAddress)
		unaryInterceptors = append(unaryInterceptors, serverinterceptor.UnaryOpentracingInterceptor())
		streamInterceptors = append(streamInterceptors, serverinterceptor.StreamOpentracingInterceptor())
	}

	grpcOptions := []grpc.ServerOption{}
	grpcOptions = append(grpcOptions,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...))
	grpcOptions = append(grpcOptions, grpc.MaxRecvMsgSize(options.maxRecvMsgSize))
	grpcOptions = append(grpcOptions, grpc.MaxSendMsgSize(options.maxSendMsgSize))
	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(options.kasp))
	if options.enableKeepAlivePolicy {
		grpcOptions = append(grpcOptions, grpc.KeepaliveEnforcementPolicy(options.kaep))
	}
	if options.creds != nil {
		grpcOptions = append(grpcOptions, grpc.Creds(options.creds))
	}
	grpcSvr := grpc.NewServer(grpcOptions...)

	app.Regist(grpcSvr)
	grpc_health_v1.RegisterHealthServer(grpcSvr, health)

	service.opts = options
	service.app = app
	service.grpcServer = grpcSvr
	service.health = health
	service.stopCh = signalutil.SetupSignalHandler()
	return service
}

func (s *service) Name() string {
	return s.opts.serverName
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Run() error {
	err := s.PreRun()
	if err != nil {
		return err
	}

	shutdown := make(chan struct{})
	serverDone, err := s.NonBlockingRun(shutdown)
	if err != nil {
		return err
	}

	<-s.stopCh
	close(shutdown)

	err = s.PreShutdown()
	if err != nil {
		return err
	}
	s.Shutdown()
	<-serverDone
	return nil
}

func (s *service) PreRun() error {
	for _, fn := range s.opts.preRunHooks {
		err := fn()
		if err != nil {
			return err
		}
	}
	s.health.SetServingStatus(s.opts.serverName, grpc_health_v1.HealthCheckResponse_SERVING)
	return nil
}

func (s *service) NonBlockingRun(shutdown <-chan struct{}) (<-chan struct{}, error) {
	lis, err := net.Listen("tcp", s.opts.listenOn)
	if err != nil {
		return nil, err
	}
	serverShutdown := make(chan struct{})
	go func() {
		defer close(serverShutdown)
		err := s.grpcServer.Serve(lis)
		select {
		case <-shutdown:
			logrus.Info("Catch signal, server will shutdown")
		default:
			panic(fmt.Sprintf("grpc serve failed, err: %v", err))
		}
	}()
	return serverShutdown, nil
}

func (s *service) PreShutdown() error {
	for _, fn := range s.opts.preShutdownHooks {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Shutdown() {
	s.grpcServer.Stop() //GracefulStop will not close imediately with stream conn
}
