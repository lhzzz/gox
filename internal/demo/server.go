package demo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"singer.com/basic/errorx"
	"singer.com/basic/meta"
	pb "singer.com/example/micro/demo"
)

type app struct {
}

func NewDemoApp() app {
	return app{}
}

func (a *app) Regist(s *grpc.Server) {
	pb.RegisterDemoServiceServer(s, a)
}

func (a *app) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	logrus.WithField("requestId", meta.GetRequestId(ctx)).Info("req:", req)
	resp := &pb.Response{}
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
