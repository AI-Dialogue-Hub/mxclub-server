package util

import (
	"mxclub/pkg/utils"
	"testing"
)

func TestGenExecutorId(t *testing.T) {
	cnt := utils.ParseInt(generateOrderID())
	t.Logf("%v", cnt)
}
