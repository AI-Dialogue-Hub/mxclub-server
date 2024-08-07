package wxpay

import (
	"fmt"
	"mxclub/pkg/utils"
	"testing"
)

func TestRefundDTO(t *testing.T) {
	val := getTestEncryptWxpayCallBackInfo(t)
	t.Logf("%v\n", val)
	jsonStr := utils.ObjToJsonStr(val)
	fmt.Printf("%v", jsonStr)
}

func getTestEncryptWxpayCallBackInfo(t *testing.T) *WxPayCallBackEncryptDTO {
	str := `{"summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"az7mOM5WFBctntP69770Fa9cQWk7RXYi/2pAD6U5H9Yo21H433y/SXvTm6QuQtPdXVgQX+WekFuj/bgIC6XAdptHtvirFdRHOidMwaGLIEucLVPOwrHOBD09nLty8KLUxNZTa1dtlxBdjUcQpmOqIpiDODAWff1k0HYNQeBtmGdZlV82pMzy+aewvi8dVk7oybLpcBrApmNGkkIkPY+zkHpdAMyId+svyXP7eOK+pO9ezakbhzFaLp0wBhma/uHYU+qAP2qhLBIpK7JnPiGXioBa4HlBalND4/I4cqZV5zPCH7K94Wo1CUDeza62ZYO5f0JpK4aKE+EhbkVhZwD8xy3DZXWmSWaAxXNhHK3lewMbCUt1M0fE6sDveHadZQxx55MKDvSSK3lGn47dcXliuYYoCA0k4ZscVA3fRM6oo087dXvQmHDljTAgSx/VFqlufexfrBKJbge5s+QF3WFyuaw3t6JTmdniJ66ylPH3T14/dGun7087NCYDFalPXmMcIoO1Pci8FJU9FpTQLWvIvlO+XiWm6N5rpiagAmSjwR6Mw9ePpxW5B7yL3BL/qfGiGXSLBQpCCr23ClZCbVpmm/bGROE=","nonce":"sLM5jAokRfXe","original_type":"transaction"},"id":"7f3344b5-3148-54b6-9146-418a83898399","create_time":"2024-08-07T03:31:37+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
