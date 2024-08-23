package utils

import (
	"fmt"
	"testing"
	"time"
)

// TestFormatDuration 测试 FormatDuration 函数
func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		duration time.Duration
		expected string
	}{
		{time.Minute * 59, "59分钟"},
		{time.Minute * 61, "1小时1分钟"},
		{time.Hour*2 + time.Minute*42 + time.Second*43, "2小时42分钟"},
		{time.Hour*24 + time.Minute*15, "1天15分钟"},
		{time.Hour*25 + time.Minute*30, "1天1小时30分钟"},
		{time.Hour*48 + time.Minute*59, "2天59分钟"},
		{time.Hour*72 + time.Minute*30, "3天30分钟"},
		{time.Hour*168 + time.Minute*45, "7天45分钟"},
		{time.Hour * 0, "0分钟"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.duration), func(t *testing.T) {
			actual := FormatDuration(tc.duration)
			if actual != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, actual)
			}
		})
	}
}
