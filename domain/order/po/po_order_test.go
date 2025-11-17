package po

import (
	"fmt"
	"testing"
)

// 测试 ExtractID 函数
func TestExtractID(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"【Id:10913696132,角色:玉米炒骨灰】", "10913696132", false},
		{"Id:10913696132,角色:玉米炒骨灰", "10913696132", false},
		{"Id:8922328330,角色:8922328330", "8922328330", false},
		{"【Id:,角色:玉米炒骨灰】", "", true},                        // 错误情况：没有 ID
		{"无有效信息", "", true},                                 // 错误情况：格式不对
		{"Id:11444852481,角色:风无痕0344", "11444852481", false}, // 错误情况：格式不对
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input=%s", tc.input), func(t *testing.T) {
			got, err := ExtractID(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractID() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.expected {
				t.Errorf("ExtractID() = %v, expected %v", got, tc.expected)
			}
		})
	}
}

// 测试 ExtractRole 函数
func TestExtractRole(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"【Id:10913696132,角色:玉米炒骨灰】", "玉米炒骨灰", false},
		{"Id:8922328330,角色:8922328330", "8922328330", false},
		{"【Id:10913696132,角色:】", "", true}, // 错误情况：没有角色名
		{"无有效信息", "", true},                // 错误情况：格式不对
		{"Id:11444852481,角色:风无痕0344", "风无痕0344", false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input=%s", tc.input), func(t *testing.T) {
			got, err := ExtractRole(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractRole() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.expected {
				t.Errorf("ExtractRole() = %v, expected %v", got, tc.expected)
			}
		})
	}
}
