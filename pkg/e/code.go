package e

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	//成员错误
	ErrorExistUser      = 10002 //用户存在
	ErrorNotExistUser   = 10003 //用户不存在
	ErrorFailEncryption = 10006 //用户密码加密错误
	ErrorNotCompare     = 10007 //密码匹配错误
	ErrorCreateUser     = 10008

	ErrorAuthCheckTokenFail    = 30001 //token 错误
	ErrorAuthCheckTokenTimeout = 30002 //token 过期
	ErrorAuthToken             = 30003
	ErrorAuth                  = 30004

	ErrorDatabase    = 40001
	ErrorGetUserInfo = 40002
	ErrorCreateTask  = 40003
)
