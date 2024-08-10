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
	str := `{"summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"C9ZalMqbXPufxWdFSqyDCCTeYDMK4x42x0SePJ874CeDSbj8a717LtJeZZklgqxfW9Z50FTcPgXGZmVppHqESxkD/tfiPRYlAbfajCrA0kvWajG5ZrMXGDusX4abYChQuz98i8jJvK/bVtMiR+S0RJbSXc7MY05WHo3B3VTs7MP0bCcgPS1SpfG/pVYDmOOKQdun1E4MqIbl2MNAdjPdwwEtG7SG+C3X73ge9DlspFHgN9egB64tN/PY3BzJ1QKa2C7LeAG/E12G6TrYJr4h1QqirAgJTuL5eVea1EkH8WC9rAewyynPbNbui0WcK99jekKEpsEyzl/w3lDEtAZM8AGI7mUIXZClbQLf/yKnDwYpR6iwk9H7zgwRMsZi2ViBB+aHjZyolIZvqstCiKFqKyrV4GFmBwrm6rzGXTQgiNX1W68VHulldj5zPpg8vwgPiIVwQ472qpC1xZ3jlC8/F/aNeOYeikmPlwdzkcffqmQu9kspgVDG6AQdaOvRHZuDMc4DrHXQTgcdS7cunYrAF7x16UlSiRGeCd99jHCS93QJPBbaENi9ZWZjC7icQ62bSdmYprqnHWWx0y1qyMQEKwxHGHrGk2jiWspm","nonce":"fv1XisQ8zfAn","original_type":"transaction"},"id":"4a5a469e-f34e-5b08-abd7-1de5b0109017","create_time":"2024-08-09T01:02:35+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
