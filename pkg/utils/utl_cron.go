package utils

import "time"

func SetTimeOut(delayTime time.Duration, fn func()) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		RecoverByPrefixNoCtx("SetTimeOut")
		select {
		case <-time.After(delayTime):
			fn()
		case <-done:

		}
	}()
	return done
}
