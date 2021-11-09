package controllers

type ResCode int64
//定义返回状态
const (
	CodeSuccess  ResCode = 1000+iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken

)
//定义返回信息
var codeMsgMap = map[ResCode]string{
	CodeSuccess: "success",
	CodeInvalidParam: "请求参数有误",
	CodeUserExist: "用户存在",
	CodeUserNotExist: "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy: "服务繁忙",
	CodeInvalidToken: "无效Token",
	CodeNeedLogin: "需要登录",
}

func (c ResCode) Msg() string{
	msg,ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
