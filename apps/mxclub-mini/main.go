package main

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"mxclub/apps/mxclub-mini/config"
	_ "mxclub/apps/mxclub-mini/controller"
	_ "mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/common/xjet"
)

func main() {
	xlog.SetOutputLevel(xlog.Ldebug)
	jet.AddMiddleware(
		xjet.CorsMiddleware,
		jet.TraceJetMiddleware,
		jet.RecoverJetMiddleware,
	)
	jet.SetFastHttpServer(&fasthttp.Server{
		Handler:            router.ServeHTTP,
		MaxRequestBodySize: config.GetConfig().File.MaxRequestBodySize * 1024 * 1024,
	})
	jet.Run(config.GetConfig().Server.Port)
}
