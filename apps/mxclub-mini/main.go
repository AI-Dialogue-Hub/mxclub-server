package main

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/apps/mxclub-mini/config"
	_ "mxclub/apps/mxclub-mini/controller"
)

func main() {
	xlog.SetOutputLevel(xlog.Ldebug)
	jet.AddMiddleware(jet.TraceJetMiddleware, jet.RecoverJetMiddleware)
	jet.Run(config.GetConfig().Server.Port)
}
