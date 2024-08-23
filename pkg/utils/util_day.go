package utils

import (
	"fmt"
	"time"
)

// FormatDuration 将一个 time.Duration 对象格式化为 "xx天xx小时xx分钟" 的形式。
// 如果总分钟数小于60，则直接返回分钟数；否则，计算天数、小时数和剩余的分钟数。
func FormatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var formatted string
	if days > 0 {
		formatted = fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		formatted += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 {
		formatted += fmt.Sprintf("%d分钟", minutes)
	}

	// 如果 formatted 为空，则说明天数、小时数和分钟数都是0，返回 "0分钟"
	if formatted == "" {
		formatted = "0分钟"
	}

	return formatted
}
