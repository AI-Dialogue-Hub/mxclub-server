package main

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/apps/mxclub-mini/config"
	_ "mxclub/apps/mxclub-mini/controller"
	"mxclub/apps/mxclub-mini/middleware"
	_ "mxclub/apps/mxclub-mini/service"
)

func main() {
	xlog.SetOutputLevel(xlog.Ldebug)
	jet.AddMiddleware(
		middleware.CorsMiddleware,
		jet.TraceJetMiddleware,
		jet.RecoverJetMiddleware,
	)
	jet.Run(config.GetConfig().Server.Port)
}
