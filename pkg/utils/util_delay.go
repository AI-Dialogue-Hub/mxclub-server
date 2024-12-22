package utils

import "time"

func DelayFunc(f func(), duration time.Duration) {
	// 使用匿名函数和 goroutine 来等待，而不会阻塞主程序
	go func() {
		defer RecoverByPrefixNoCtx("DelayFunc")
		// 等待指定的持续时间
		<-time.After(duration)
		// 执行传入的函数
		f()
	}()
}
