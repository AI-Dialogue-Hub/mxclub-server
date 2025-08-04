package xjet

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
)

func CorsMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return jet.JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		handleCors(ctx)
		// 如果是OPTIONS请求，直接返回
		if string(ctx.Method()) == "OPTIONS" {
			ctx.Response.Header.SetServer("JetServer")
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		next.ServeHTTP(ctx)
	}), nil
}

func handleCors(ctx *fasthttp.RequestCtx) {
	// 设置允许跨域请求的响应头
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, token, Authorization")
	ctx.Response.Header.Set("Access-Control-Max-Age", "86400")
}
