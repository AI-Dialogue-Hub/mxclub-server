package util

import (
	"mxclub/pkg/utils"
	"testing"
	"time"
)

func TestGenExecutorId(t *testing.T) {
	cnt := utils.ParseInt(generateOrderID())
	t.Logf("%v", cnt)
}

func TestTime(t *testing.T) {
	now := time.Now()
	time.Sleep(time.Microsecond * 5)
	t.Logf("%v", time.Since(now))
}
