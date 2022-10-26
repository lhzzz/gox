package serverinterceptor

import (
	"context"

	"singer.com/basic/errorx"
	"singer.com/basic/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//error拦截器可以自动打印错误日志，并且将自定义类型错误对应的脱敏信息返回给调用者 (自定义错误类型为errorx)
func UnaryErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		log.RpcErrorf(ctx, "%+v", err)
		causeErr := errorx.Cause(err)                  // err类型
		if e, ok := causeErr.(*errorx.CodeError); ok { //自定义错误类型
			//转成grpc err
			err = status.Error(codes.Code(e.GetErrCode()), e.GetUsrMsg())
			e.SetTrailer(ctx)
		}
	}
	return resp, err
}

func StreamErrorInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	if err != nil {
		log.RpcErrorf(ss.Context(), "%+v", err)
		causeErr := errorx.Cause(err)                  // err类型
		if e, ok := causeErr.(*errorx.CodeError); ok { //自定义错误类型
			//转成grpc err
			err = status.Error(codes.Code(e.GetErrCode()), e.GetUsrMsg())
			e.SetTrailer(ss.Context())
		}
	}
	return err
}
