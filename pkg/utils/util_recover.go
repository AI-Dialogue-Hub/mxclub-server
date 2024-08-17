package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
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
