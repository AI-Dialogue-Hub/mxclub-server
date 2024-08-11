package main

import (
	jetContext "github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/config"
	_ "mxclub/apps/mxclub-admin/controller"
	"mxclub/apps/mxclub-admin/middleware"
	_ "mxclub/apps/mxclub-admin/service"
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
	jet.AddPostJetCtxInitHook(func(ctx jetContext.Ctx) {
		xmysql.SetLoggerPrefix(ctx.Logger().ReqId)
	})
	jet.Run(config.GetConfig().Server.Port)
}
