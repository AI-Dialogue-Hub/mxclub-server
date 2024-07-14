package main

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/config"
	_ "mxclub/apps/mxclub-admin/controller"
	"mxclub/apps/mxclub-admin/middleware"
	_ "mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/common/xjet"
)

func main() {
	jet.AddMiddleware(
		xjet.CorsMiddleware,
		middleware.AuthMiddleware,
		jet.TraceJetMiddleware,
		jet.RecoverJetMiddleware,
	)
	jet.Run(config.GetConfig().Server.Port)
}
