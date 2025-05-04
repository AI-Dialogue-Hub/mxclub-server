// Package middleware 记录后台人员的操作
package middleware

import (
	"encoding/json"
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"mxclub/apps/mxclub-admin/config"
	"mxclub/domain/user/po"
	"strings"
)

func OperatorMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return jet.JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		if err := handleOperator(ctx); err == nil {
			next.ServeHTTP(ctx)
		} else {
			ctx.Response.Header.SetServer("JetServer")
			ctx.Response.Header.Set("Content-Type", constant.MIMETextPlain)
			ctx.SetBodyString(err.Error())
		}
	}), nil
}

func handleOperator(ctx *fasthttp.RequestCtx) (err error) {
	path := string(ctx.Path())

	// 判断是否为开放API或上传接口，如果是则只打印基本信息
	if config.IsOpenApi(path) {
		xlog.Info("[OperatorMiddleware#handleOperator]",
			"requestURI", string(ctx.Request.URI().Path()),
			"method", string(ctx.Method()))
		return nil
	}

	// 假设非上传请求通过方法类型判断，可根据实际情况调整
	if string(ctx.Method()) == "POST" && strings.Contains(path, "/upload") {
		// 对于上传请求，不做额外处理
		return nil
	} else {
		// 打印入参，美化输出
		args := ctx.QueryArgs()
		params := make(map[string]string)
		args.VisitAll(func(key, value []byte) {
			params[string(key)] = string(value)
		})

		prettyParams, errMarshal := json.MarshalIndent(params, "", "  ")
		if errMarshal != nil {
			xlog.Error("Failed to marshal request parameters")
		}

		userPO, err := FetchUserInfoByCtx(ctx)

		if err != nil {
			return err
		}

		// 记录用户操作了哪个接口以及请求参数
		xlog.Infof("[OperatorMiddleware#handleOperator] User %s (%d) accessed %s with method %s\nRequest Params:\n%s",
			userPO.Name, userPO.ID,
			string(ctx.Request.URI().Path()),
			string(ctx.Method()),
			string(prettyParams))
	}
	return nil
}

func FetchUserInfoByCtx(ctx *fasthttp.RequestCtx) (*po.User, error) {
	// 获取用户信息
	rawUserInfo := ctx.UserValue("tokenInfo")
	if rawUserInfo == nil {
		xlog.Error("[OperatorMiddleware#handleOperator] userInfo is invalid")
		return nil, errors.New("未登录")
	}

	// 类型断言获取用户详细信息
	userInfo := rawUserInfo.(*AuthToken)
	userPO := userInfo.UserPO
	return userPO, nil
}
