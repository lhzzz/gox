package errorx

var message map[ErrorCode]string

func init() {
	message = make(map[ErrorCode]string)
	message[OK] = "SUCCESS"
	message[SERVER_COMMON_ERROR] = "服务器开小差啦, 稍后再来试一试"
	message[REUQEST_PARAM_ERROR] = "参数错误"
	message[TOKEN_EXPIRE_ERROR] = "token失效, 请重新登陆"
	message[TOKEN_GENERATE_ERROR] = "生成token失败"
	message[DB_ERROR] = "数据库繁忙, 请稍后再试"
}

func MapErrMsg(errcode ErrorCode) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦, 稍后再来试一试"
	}
}

func IsCodeErr(errcode ErrorCode) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
