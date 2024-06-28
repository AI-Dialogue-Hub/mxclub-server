package api

import (
	"net/http"
)

const (
	// 4XX
	ErrorCodeBadRequest   = http.StatusBadRequest   // 400
	ErrorCodeUnauthorized = http.StatusUnauthorized // 401
	ErrorCodeForbidden    = http.StatusForbidden    // 403
	ErrorCodeNotFound     = http.StatusNotFound     // 404

	// 5XX
	ErrorCodeInternal = http.StatusInternalServerError // 5XX

	// User
	ErrorCodeUser         = 10001 // 用户模块错误码，
	ErrorCodeSSOFail      = 10002 // SSO 鉴权失败，无效的token
	ErrorCodeNoUser       = 10003 // 获取用户信息失败
	ErrorCodeTokenExpired = 10004 // Token 已经过期

	// App
	ErrorCodeAppNameExisted    = 20001 // 应用名已经存在
	ErrorCodeAppPackageExisted = 20002 // 应用的包名已经存在
	ErrorCodeAppBundleExisted  = 20003 // 应用的bundle 已经存在
	ErrorCodeAppCantEdit       = 20004 // 当前状态无法编辑
	ErrorCodeAppEditPlatform   = 20005 // 平台不能修改

	// 一键登录
	ErrorCodeInvalidAppStatus = 30001 // APP 状态不可用
	ErrorCodeNoRsaPublicKey   = 30002 // 未配置RSA 公钥
	ErrorCodeRemoteError      = 30003 // 请求外部系统发生异常
	ErrorCodeRemoteFail       = 30004 // 外部系统返回失败状态
	ErrorCodeRecordNotFound   = 30005 // RecordNotFound
)

func ErrorBadRequest(reqId string, msgs ...string) *Response {
	msg := "Bad request"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return Error(reqId, ErrorCodeBadRequest, msg)
}

func ErrorUnauthorized(reqId string, msgs ...string) *Response {
	msg := "Unauthorized"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return Error(reqId, ErrorCodeUnauthorized, msg)
}

func ErrorForbidden(reqId string, msgs ...string) *Response {
	msg := "您没有权限，请联系管理员开通."
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return Error(reqId, ErrorCodeForbidden, msg)
}

func ErrorNotFound(reqId string, msgs ...string) *Response {
	msg := "Not found"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return Error(reqId, ErrorCodeNotFound, msg)
}

func ErrorInternal(reqId string, msgs ...string) *Response {
	msg := "Internal error"

	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return Error(reqId, ErrorCodeInternal, msg)
}
