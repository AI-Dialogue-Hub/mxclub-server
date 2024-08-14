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
	str := `{"create_time":"2024-08-14T14:22:31+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"vmRCco8vMwKuQKuUcJgkBfNwtqAcDJEDX924d7I2L2zTs2TGZkhBIMPusZZ9dDXZAlyaCRjdlLaqJw+24MtvD3evdaP/g+fSrRADd6uViGafvnQh52pJ/fwhbheQH5a5qF46jOn9EylttezekkWVR0KjKnTRmYE/ReXdYkJfHk30CHVixgeRTVcY3eQcqgEamejzfv8otqEpiCSRtfWlmsv9VBDLZZdk1/6iqxRQvnRNJZp+vBvM4JlO3xKZd7l8hyU1z7VFe5oxYnf5rcltoX/tSLpV0V12fNH/tAJ6OyKC5X/YuHAD9PZAVh2/pwKXUsGUSVZYUHUnFAatdOadmbwL5p3vjUfrveNI7gf1Br9/PepFo+hsQZIICp4ccfP2wR5OmoUO7LDbPtWuO6RJevAarmjXsMMJYaaq0ZaTtp/J27uQh7J7mPGWUT88lM7HSyIPV/u+OeCAdULwdN1GR3hQgJ6l7hL6Gp1dpY8C5uDXVoZqyC2FTsEkU4H4g87gcU5g23qLZOKdIgO6hsn2bYzEiOVE/RLwBtdoDHXTw18hTZkuWGkSEFb+Rf47zO2vgOt4fC0gmsrAvPof7we1zxu5UqIUECO/Pg==","nonce":"egZWZ2tozRrT","original_type":"transaction"},"id":"f963bec1-d826-535c-ac18-74d24db7c89c"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
