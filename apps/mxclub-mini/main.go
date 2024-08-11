package main

import (
	jetContext "github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
	"mxclub/apps/mxclub-mini/config"
	_ "mxclub/apps/mxclub-mini/controller"
	"mxclub/apps/mxclub-mini/middleware"
	_ "mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/common/xmysql"
)

func main() {
	jet.AddMiddleware(
		xjet.CorsMiddleware,
		middleware.AuthMiddleware,
		jet.TraceJetMiddleware,
		jet.RecoverJetMiddleware,
	)
	jet.SetFastHttpServer(&fasthttp.Server{
		Handler:            router.ServeHTTP,
		MaxRequestBodySize: config.GetConfig().File.MaxRequestBodySize * 1024 * 1024,
	})
	jet.AddPostJetCtxInitHook(func(ctx jetContext.Ctx) {
		xmysql.SetLoggerPrefix(ctx.Logger().ReqId)
	})
	jet.Run(config.GetConfig().Server.Port)
}
