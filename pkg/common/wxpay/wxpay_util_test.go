package wxpay

import "testing"

func TestGenerateUniqueOrderNumber(t *testing.T) {
	t.Logf("%v", generateUniqueOrderNumber())
}
