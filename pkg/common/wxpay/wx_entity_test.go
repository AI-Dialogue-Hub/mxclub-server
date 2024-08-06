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
	str := `{"event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"1ot0RoQeISENIG3jZ1EftK2YBZn0j9FxI3Obz3sty+0XTZqbPvWMMDOLEM8IUt/TpkXC/WJpwbLyK22xXAyiaC/jHPwohZpzQ1n1gRZWZRpzZcsb+6fMlu5KZZDDDjwvbF38cbiJDZvmB/gOEWXFl4nB5joHAsg+/1DEsTp5izia+flea7k+nk+K6NWdiVZE8rgXdI/i3xc3oqHdY0bFwQCVU5wiDMvl7jCngo6HVecDTapkw0mEDMJTmnSecRAgsqfm2O8xmAHWyoWf8O6BD12Ut9rtMgscVTzvYKefdknk9MKRiDnadzTfj5+2U/jOcK5N0vAMrZqRw3zAvPnWHH/DOABRe/poYFnnTzsxB40A4k93CF6T1G4aQT7723MuVZuWiTE0tzQWztwi+n3BpiIZRdeJwKqPF8S0Wd1IaTTmVqJJf8WQ2mEl/GW+7XHGzk8jG53nV3XateqDNag/DKNwv3lXIv+9TBUA9LLXd21AWZg23vTgQ1YnthOBwujMTSMSZG0BSHVXATJSWcMn44vJBZ2wMXEVxtpVuo6BTyYmzJdhvupl19Qu/8KrSb+cmB/vBYf44bLOE2+YhwjJsvIVzXU=","nonce":"RmFDhTdYA0Jr","original_type":"transaction"},"id":"d548c219-55b7-59e2-a2f9-c18059bc8429","create_time":"2024-07-19T20:12:42+08:00","resource_type":"encrypt-resource"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
