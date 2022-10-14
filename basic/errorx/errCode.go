package errorx

/**(前3位代表业务,后三位代表具体功能)**/

type ErrorCode uint32

//全局错误码
const (
	OK                   ErrorCode = 200
	SERVER_COMMON_ERROR  ErrorCode = 100001
	REUQEST_PARAM_ERROR  ErrorCode = 100002
	TOKEN_EXPIRE_ERROR   ErrorCode = 100003
	TOKEN_GENERATE_ERROR ErrorCode = 100004
	DB_ERROR             ErrorCode = 100005
)

//用户模块
