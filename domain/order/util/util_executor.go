package util

import (
	"math/rand"
	"strconv"
	"time"
)

// 生成一个唯一的打手ID
func generateDasherID() string {
	// 获取当前时间戳（秒级）
	timestamp := time.Now().Unix()

	// 生成一个3位的随机数
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(900) + 100 // 生成一个100到999之间的随机数

	// 将时间戳和随机数组合生成打手ID
	dasherID := strconv.FormatInt(timestamp, 10) + strconv.Itoa(randomNum)

	// 取打手ID的后8位，确保ID不太长且唯一
	if len(dasherID) > 8 {
		dasherID = dasherID[len(dasherID)-8:]
	}

	return dasherID
}
