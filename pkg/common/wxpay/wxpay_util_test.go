package wxpay

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

// 随机生成32位字符
func TestGenerateRandomString(t *testing.T) {
	bytes := make([]byte, 16) // 16字节 = 32个十六进制字符
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(bytes))
}
