package serverinterceptor

import (
	"context"

	"github.com/sirupsen/logrus"
	"singer.com/basic/errorx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errorx.Cause(err)                  // err类型
		if e, ok := causeErr.(*errorx.CodeError); ok { //自定义错误类型
			//转成grpc err
			err = status.Error(codes.Code(e.GetErrCode()), e.GetUsrMsg())
			e.SetTrailer(ctx)
		}
		logrus.Errorf("[RPC-ERR] %+v", err)
	}
	return resp, err
}
