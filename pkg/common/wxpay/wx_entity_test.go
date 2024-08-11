package wxpay

import (
	"fmt"
	"mxclub/pkg/utils"
	"testing"
)

func testRefundDTO(t *testing.T) {
	val := getTestEncryptWxpayCallBackInfo(t)
	t.Logf("%v\n", val)
	jsonStr := utils.ObjToJsonStr(val)
	fmt.Printf("%v", jsonStr)
}

func getTestEncryptWxpayCallBackInfo(t *testing.T) *WxPayCallBackEncryptDTO {
	str := `{"summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"zMwzl8o6CKa37X8EKRcLVEAv+1a7LpzAP7TkX3ZEEpbtyQLetKShUSSTl6xwbfAMlrZsL0PKvNaMcCZUJjjvwGLYCI4GBx64MQgOFjzd8iWLwjmFQwzVEN4KAJvPlxXq92kHN0GoDrp6Q9MNNB0eIUmVft8rrU2RsRFvn8QwVH5s9aUuuKAUFUNMYs2jHO0JvIkJz8RoBJpAmKsHPNf/L5xlG71MFNC7v9yhp1A1X2ThYCoqmbkriS1SmLT/OGhaZTcfz6VdJ79xhsZpEYwguHfGfaIJQcgmWf1lB3Dk+Y4/uNzqDq48EUt6QgrG0dxsNomD9VlHdHLCX4GBGVPxBjcDn4KwKBRW9bN5pN6kX+EHzUel9Ue0pREuRziQBvT7mMflWdi6Rp2sB/d7UTp1GrV32o2a+1U/qWxXrxOGVQIF0fYkn/PQWXD7HE5MOakRO800UfbX9/Ouoj41G0RWH4bnD5vuD0RGV1pOdv1QArvv8jWDizT2hG59dfyUwHeYE8BATeOWmVDVD1d20lQ/MqIro2OyLfylZniqajbs3fR5VcuyJW/0BWx7YzwBrTRh3/FdmmN63KaGRiX0aEs4qaXcO9AKdEc=","nonce":"sue6aFnccw2B","original_type":"transaction"},"id":"a32d1b69-9e56-57bd-8d0b-970becd4ede8","create_time":"2024-08-11T09:56:32+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
