package util

import (
	"fmt"
	"math/rand"
	"time"
)

func generateOrderID() string {
	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 生成一个3位的随机数
	rand.Seed(time.Now().UnixNano())

	randomNum := rand.Intn(10000) // 生成一个0到9999的随机数

	// 组合时间戳和随机数生成订单编号
	orderID := fmt.Sprintf("%d%04d", timestamp, randomNum)
	return orderID
}
