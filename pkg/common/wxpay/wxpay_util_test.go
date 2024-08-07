package wxpay

import (
	"mxclub/pkg/utils"
	"testing"
)

func TestGenerateUniqueOrderNumber(t *testing.T) {
	t.Logf("%v", utils.ParseInt(GenerateUniqueOrderNumber()))
}

func TestGenerateOutRefundNo(t *testing.T) {
	t.Logf("%v", GenerateOutRefundNo())
}
