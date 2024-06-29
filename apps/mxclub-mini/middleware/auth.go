package middleware

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"mxclub/pkg/api"
)

const AuthHeaderKey = "Authorization"

func AuthMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return jet.JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		if err := handleJwtAuth(ctx); err == nil {
			next.ServeHTTP(ctx)
		} else {
			ctx.Response.Header.SetServer("JetServer")
			ctx.Response.Header.Set("Content-Type", constant.MIMEApplicationJSON)
			ctx.SetBodyString(err.Error())
		}
	}), nil
}

func handleJwtAuth(ctx *fasthttp.RequestCtx) (err error) {
	authInfo := string(ctx.Request.Header.Peek(AuthHeaderKey))
	logger := xlog.NewWith("auth_middleware")
	if authInfo == "" {
		logger.Error("empty Authorization")
		err = api.ErrorUnauthorized(logger.ReqId)
	}

	return
}
