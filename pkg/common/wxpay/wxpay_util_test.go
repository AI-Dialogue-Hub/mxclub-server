package wxpay

import (
	"github.com/fengyuan-liang/GoKit/collection/sets"
	"mxclub/pkg/utils"
	"testing"
)

func TestGenerateUniqueOrderNumber(t *testing.T) {
	cnt := 10
	set := sets.NewHashSet[int]()
	for i := 0; i < cnt; i++ {
		parseInt := utils.ParseInt(GenerateUniqueOrderNumber())
		if set.Contains(parseInt) {
			t.Fatalf("durable number: %v, set is: %v", parseInt, set)
		} else {
			set.Add(parseInt)
		}
	}
}

func TestGenerateOutRefundNo(t *testing.T) {
	t.Logf("%v", GenerateOutRefundNo())
}
