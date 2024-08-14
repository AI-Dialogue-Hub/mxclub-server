package enum

import "testing"

func TestDeductStatus_DisPlayName(t *testing.T) {
	rule, exists := DiscountRules[""]
	t.Logf("rule:%+v, exists:%v", rule, exists)
}
