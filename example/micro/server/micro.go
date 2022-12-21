package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"singer.com/basic/errorx"
	"singer.com/basic/meta"
	"singer.com/example/micro/demo"
	"singer.com/rpc/micro"
)

type app struct {
}

func (a *app) Regist(s *grpc.Server) {
	demo.RegisterDemoServiceServer(s, a)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	a := app{}
	s := micro.NewService(&a, micro.Name("demo"),
		micro.EnableMetric("demo"),
		micro.SlowThreshold(500*time.Millisecond),
		micro.Timeout(3*time.Second),
		micro.SetMaxConnectionIdle(5*time.Minute),
		micro.SetMaxRecvMsgSize(128),
		micro.SetMaxSendMsgSize(64),
		micro.PreRunHooks(func() error {
			logrus.Info("pre running, what to do")
			return nil
		}),
		micro.PreShutdownHooks(func() error {
			logrus.Info("pre shutdown, what to finish")
			time.Sleep(5 * time.Second)
			return nil
		}),
	)
	s.Run()
}

func (a *app) Get(ctx context.Context, req *demo.Request) (*demo.Response, error) {
	logrus.WithField("requestId", meta.GetRequestId(ctx)).Info("req:", req)
	resp := &demo.Response{}
	if req.Msg == "slow" {
		time.Sleep(1 * time.Second)
		resp.Data = "slow-rsp"
	} else if req.Msg == "timeout" {
		time.Sleep(5 * time.Second)
	} else if req.Msg == "maxSend" {
		str := "1234567898765432" //16
		str += str                //32
		str += str                //64
		resp.Data = str
	} else if req.Msg == "error" {
		return resp, errorx.NewErrCode(errorx.DB_ERROR)
	} else if req.Msg == "error-msg" {
		return resp, errorx.NewErrCodeMsg(errorx.SERVER_COMMON_ERROR, "server internal error")
	} else if req.Msg == "wrap-error" {
		return resp, errorx.Wrap(errorx.REUQEST_PARAM_ERROR, "param wrap failed")
	} else if req.Msg == "wrapf-error" {
		return resp, errorx.Wrapf(errorx.REUQEST_PARAM_ERROR, "param wrapf failed")
	} else if req.Msg == "panic" {
		panic("msg call panic")
	}
	return resp, nil
}
