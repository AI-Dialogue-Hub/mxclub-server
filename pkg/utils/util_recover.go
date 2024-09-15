package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"runtime/debug"
)

func RecoverAndLogError(ctx jet.Ctx) {
	if r := recover(); r != nil {
		ctx.Logger().Error("Recovered from panic:", r)
		debug.PrintStack()
	}
}

func RecoverWithPrefix(ctx jet.Ctx, prefixInfo string) {
	if r := recover(); r != nil {
		ctx.Logger().Errorf("[%v]Recovered from panic: %v", prefixInfo, r)
		debug.PrintStack()
	}
}

func RecoverByPrefix(logger *xlog.Logger, prefixInfo string) {
	if r := recover(); r != nil {
		logger.Errorf("[%v]Recovered from panic: %v", prefixInfo, r)
		debug.PrintStack()
	}
}

func RecoverByPrefixNoCtx(prefixInfo string) {
	if r := recover(); r != nil {
		xlog.Errorf("[%v]Recovered from panic: %v", prefixInfo, r)
		debug.PrintStack()
	}
}
