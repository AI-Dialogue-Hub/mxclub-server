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
	str := `{"create_time":"2024-08-05T14:34:39+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"M6cSDm4j7BxqLs4uI5AJ/OG/Tkmqnq0sDKlxN7Pep8uRySKAdfbtSykWIc3c7RZ8VtifkImpzWHKHCGgCiML00wLKjWH6y/R47/8gQNe7MSzISL3shKdWMocj6M4c3fhmHbMaQyGm20zj8RHmsMDp06hs9xNiTnPVEXuQ6X1r23J3jcyJbyn2YSSTujz236DgQXg7AuABDst0eQhUfK9I/wQT6arCgfXABTLTHMOFQ7TemvUIKYrI6MRV9nT2UbsngLXAIG891EtB86QdDTiWWnX4ff/cYpWLUiQzmwf3LLJHyI950t18CuFv8VHQaiUGsK3JO6KuPbwfLWoECL77SgyteIFV8I8p9SlFCQrX6B66dyEnh89ITNstctLMFecv30K+vRFkBM261SBdHCiyAq2EXAbJuUprJK3QjB/IS6c/7DSI8/SoCF+YbVmvac8rIieGVDjRd9EU4QTz3suxIi9kunX4YhtjWcm7K/D3196+i8NFHHYlAJi/ZjvEtyvC13MublH43md5LhGiph7b5QARex8bGdW4lqpEgg/ugTSj9C5W2mxZqu1gnh06WiqvnwQZ3KkJYptywp3UUTCg0EnDM5KVbY=","nonce":"hS7Y8RLMelZq","original_type":"transaction"},"id":"9150f1e1-b028-5019-b323-dddb6579996a"}`
	val, err := utils.JsonStrToObj[WxPayCallBackEncryptDTO](str)
	if err != nil {
		t.Logf("error, %v", err)
		t.Failed()
	}
	return val
}
