package middleware

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
)

func CorsMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return jet.JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		handleCors(ctx)
		// 判断请求方法是否为 OPTIONS，如果是，则直接返回，不进行后续处理
		if string(ctx.Method()) == "OPTIONS" {
			return
		}
		next.ServeHTTP(ctx)
	}), nil
}

func handleCors(ctx *fasthttp.RequestCtx) {
	// 设置允许跨域请求的响应头
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform")
}
