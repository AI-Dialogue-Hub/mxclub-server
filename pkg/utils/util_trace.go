package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"time"
)

func TraceElapsed(ctx jet.Ctx, operation string) func() {
	if ctx == nil || ctx.Logger() == nil {
		return func() {} // 返回空函数避免panic
	}

	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		ctx.Logger().Infof("[%s] elapsed time: %v", operation, elapsed)
	}
}

func TraceElapsedWithPrefix(logger *xlog.Logger, operation string) func() {
	if logger == nil {
		return func() {} // 返回空函数避免panic
	}

	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		logger.Infof("[%s] elapsed time: %v", operation, elapsed)
	}
}
