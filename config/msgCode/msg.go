package msgCode

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	UN_LOGGING:                     "用户未登录",
	ERROR_NOT_EXIST_AUTH:           "账号或密码不存在",
	ERROR_AUTH_TOKEN_CHECK_FAIL:    "token校验失败",
	ERROR_AUTH_TOKEN_CHECK_TIMEOUT: "token过期",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_EXIST_TOKEN:              "已存在该token",
	ERROR_NOT_EXIST_TOKEN:          "token不存在",
	ERROR_NOT_EXIST_TOKENID:        "tokenId不存在",
	ERROR_SYNC_YAPI:                "Yapi同步接口出错",
	ERROR_SYNC_RUNNING:             "Yapi接口同步中...",
	ERROR_NOT_EXIST_GROUP:          "业务线不存在",
	ERROR_EXIST_AUTH:               "用户已存在",
	ERROR_EXIST_TAG:                "标签已存在",
	ERROR_NOT_EXIST_TAG:            "标签不存在",
	ERROR_NOT_EXIST_ARTICLE:        "文章不存在",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
