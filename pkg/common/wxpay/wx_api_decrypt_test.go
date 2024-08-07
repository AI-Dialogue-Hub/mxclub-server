package wxpay

import (
	"bytes"
	"context"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"io"
	"log"
	"mxclub/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testConfigDTO struct {
	Config *WxPayConfig `yaml:"wx_pay_config"`
}

var (
	testConfigPath = "E:\\workspace\\goland\\config\\mxclub.yml"
	testConfig     = new(testConfigDTO)
)

func setUp() {
	// 读取配置文件
	if err := utils.YamlToStruct(testConfigPath, testConfig); err != nil {
		log.Fatalf("config parse error:%v", err.Error())
	}
	InitWxPay(testConfig.Config)
}

func testDecryptData(t *testing.T) {
	setUp()
	info := getTestEncryptWxpayCallBackInfo(t)
	handler := NewWxPayCertHandler(testConfig.Config)
	req := httptest.NewRequest(
		http.MethodGet, "http://127.0.0.1", io.NopCloser(bytes.NewBuffer(utils.MustObjToByte(info))),
	)
	req.Header.Add("test", "test")
	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), req, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		fmt.Println(err)
		return
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(*transaction)
	fmt.Println(utils.ObjToJsonStr(transaction.TransactionId))
	fmt.Println(utils.ObjToMap(*transaction))
	fmt.Println(utils.MustMapToObj[payments.Transaction](utils.ObjToMap(*transaction)))
}
