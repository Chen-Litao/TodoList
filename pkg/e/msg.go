package e

var CodeMsg = map[int]string{
	SUCCESS:       "操作成功",
	ERROR:         "操作失败",
	InvalidParams: "参数无效",

	ErrorExistUser:             "用户已存在",
	ErrorNotExistUser:          "用户不存在",
	ErrorCreateUser:            "用户创建失败",
	ErrorFailEncryption:        "密码加密无效",
	ErrorNotCompare:            "对比无效",
	ErrorAuthCheckTokenFail:    "token校验无效",
	ErrorAuthCheckTokenTimeout: "token过期",
	ErrorAuthToken:             "授权token失效",
	ErrorAuth:                  "Token错误",
	ErrorDatabase:              "数据库操作出错,请重试",
	ErrorGetUserInfo:           "获取用户信息失败",
	ErrorCreateTask:            "创建任务失败",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := CodeMsg[code]
	if ok {
		return msg
	}
	return CodeMsg[ERROR]
}
