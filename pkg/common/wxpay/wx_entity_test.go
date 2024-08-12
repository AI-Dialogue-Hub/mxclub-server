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
	str := `{"create_time":"2024-08-12T23:23:00+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"mNGAosSSTGPGAxvegplmKCUQ5JLt6X4xtwY2uK7irFpMp5hVyOEjmYH3UpUqJFnuMufCsZ1uMs8vgF1rjmtqFndFh0DOV2dzEPN4dDg57Rxce/jHV9hM9Sskm/AvNcjrmGLLwh10Z0I9RSHgCwPW/Vc4PBrQV5sTYjeF7HBblcb/qzI3NbfEXOaBqwZxskRKXtDyatDGEFwWYIR4bTkcac+xaP+91JOth1mT76yRkjGHHe3kBKu4fuai2WcQe6Gs587Y/PrhM/GYksIybWAg2ikl9sknAqmAqpI491dFACtI0y7BTx7NVCkloG29w7TBpjnpdnZERQ3CXSFmx7Jb/hQF9ucfRr3aHwnJ61C0yHWO8rtBCqvpNpQs4fCOeBI7kGrZxJ1Sf9SvVcde4Jrow/8O47CEy/c8vWFn5Hffhmwb0nKIW0FtGx4E8Em4phCr52ekL8bD0CS9DNK67802ZUbz1nXAo54dUYNP2UkZ5j36o2PLYhzVRmNTSMrqolvF7h9PXNP8H1HJ4nWLxIhu+DP9tsNuZ7pqSgbXxcS2v7PI+u3Pt9uf2VrEtQE3T1Ob6a3iPVcEymKzQentPgrFScCadT9nYu71","nonce":"AE8IXeZmq1hU","original_type":"transaction"},"id":"1249eecb-6edd-5161-96fb-e863f978e0ba"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
