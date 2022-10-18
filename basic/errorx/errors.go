package errorx

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/**
常用通用固定错误
*/

/*
  错误处理的设计需要考虑：
1、后台尽量简洁，减少 err != nil
2、前端和后台想要看到的错误提示是不一样的,前端显示脱敏错误，后台显示内部错误
3、grpc的error返回格式
	code
	msg
	[]any 必须是proto格式

目前的设计：
	1、proto定义中不包含错误
	2、grpc 接口通过errorx.Wrapf 返回错误，通过grpc拦截器自动打印错误日志，并将错误码对应的脱敏错误返回给前端
*/

type CodeError struct {
	errCode ErrorCode //错误码
	errMsg  string    //内部错误信息
}

const (
	kErrorxTrailerKey = "errorx-message"
)

//返回给前端的错误码
func (e *CodeError) GetErrCode() ErrorCode {
	return e.errCode
}

//返回给前端的错误信息
func (e *CodeError) GetUsrMsg() string {
	return MapErrMsg(e.errCode)
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
}

func (e *CodeError) SetTrailer(ctx context.Context) {
	grpc.SetTrailer(ctx, metadata.Pairs(kErrorxTrailerKey, e.errMsg))
}

func NewErrCodeMsg(errCode ErrorCode, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}

func NewErrCode(errCode ErrorCode) *CodeError {
	return &CodeError{errCode: errCode}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: SERVER_COMMON_ERROR, errMsg: errMsg}
}

func Wrapf(ce *CodeError, format string, args ...interface{}) error {
	return errors.Wrapf(ce, format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}
