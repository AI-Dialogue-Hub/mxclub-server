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
