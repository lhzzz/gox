package client

import (
	"context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"singer.com/basic/trace"
	"singer.com/rpc/clientinterceptor"
)

type client struct {
	conn *grpc.ClientConn

	//todo: need a service resolver
	//resolver
}

type Client interface {
	Conn() *grpc.ClientConn
}

func NewClient(target string, opts ...ClientOption) (Client, error) {
	var c client
	options := c.dialOptions(opts...)
	conn, err := grpc.DialContext(context.Background(), target, options...)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return &c, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *client) dialOptions(opts ...ClientOption) []grpc.DialOption {
	opt := newDefaultClientOptions()
	for _, o := range opts {
		o(&opt)
	}

	gOptions := make([]grpc.DialOption, 0)
	gOptions = append(gOptions, grpc.WithTransportCredentials(opt.creds))
	if opt.enableTrace {
		gOptions = append(gOptions, grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(),
				otgrpc.SpanDecorator(otgrpc.SpanDecoratorFunc(trace.SpanDecoratorError)))))
	}
	if opt.block {
		gOptions = append(gOptions, grpc.WithBlock())
	}
	if opt.maxSendMsgSize > 0 {
		gOptions = append(gOptions, grpc.MaxCallSendMsgSize(opt.maxSendMsgSize))
	}
	if opt.maxRecvMsgSize > 0 {
		gOptions = append(gOptions, grpc.MaxCallRecvMsgSize(opt.maxRecvMsgSize))
	}

	unaryInterceptors := make([]grpc.UnaryClientInterceptor, 0)
	streamInterceptors := make([]grpc.StreamClientInterceptor, 0)

	if opt.enableMeta {
		unaryInterceptors = append(unaryInterceptors, clientinterceptor.UnaryMetaInterceptor)
		streamInterceptors = append(streamInterceptors, clientinterceptor.StreamMetaInterceptor)
	}

	if opt.breaker != nil && opt.accept != nil {
		unaryInterceptors = append(unaryInterceptors, clientinterceptor.BreakerInterceptor(opt.breaker, opt.accept))
	}
	if opt.slowThreshold > 0 {
		clientinterceptor.SetSlowThreshold(opt.slowThreshold)
		unaryInterceptors = append(unaryInterceptors, clientinterceptor.SlowlogInterceptor)
	}
	if opt.timeout > 0 {
		unaryInterceptors = append(unaryInterceptors, clientinterceptor.TimeoutInterceptor(opt.timeout))
	}
	gOptions = append(gOptions, grpc.WithChainUnaryInterceptor(unaryInterceptors...))
	gOptions = append(gOptions, grpc.WithChainStreamInterceptor(streamInterceptors...))
	gOptions = append(gOptions, opt.dialOptions...)
	return gOptions
}
