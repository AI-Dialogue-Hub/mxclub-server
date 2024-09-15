package wxpay

import (
	"bytes"
	"context"
	"fmt"
	jetContext "github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"io"
	"mxclub/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRefunds(t *testing.T) {
	info := getInfoV2(t)
	requestCtx := &fasthttp.RequestCtx{Request: fasthttp.Request{}, Response: fasthttp.Response{}}
	ctx := jetContext.NewContext(requestCtx, xlog.NewWith("text"))
	err := Refunds(ctx, info, GenerateOutRefundNo(), "协商一致进行退款!")
	if err != nil {
		ctx.Logger().Errorf("%v", err)
	}
}

func getInfo(t *testing.T) *payments.Transaction {
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
		return nil
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(*transaction)
	fmt.Println(utils.ObjToJsonStr(transaction.TransactionId))
	fmt.Println(utils.ObjToMap(*transaction))
	fmt.Println(utils.MustMapToObj[payments.Transaction](utils.ObjToMap(*transaction)))
	return transaction
}

func getInfoV2(t *testing.T) *payments.Transaction {
	setUp()
	//str := `{"algorithm":"AEAD_AES_256_GCM","associated_data":"transaction","ciphertext":"vaFgzv8kppgZmD51S0G9cYnf22qyDVOoaF6cyN19gYGInb8NTu0RFzbwn5UgKV/qzdsYW/HKkVnkVgT+lKB/ILmNu1+qGkCmtJCt+Cj6nnFf2aEUANJsNskrTlGzEqblxYozS12uaomt1R/9zRvgYPrS3tJxTTfaJW5mnOAgQ7lyIIIF4OlsFtdI7BXCIjfAlvJ7OsKdTB9BPU7N3943flKuHd/t7AHcqL9SOe1srlWESLLbn+yoWnDjVYdWY2XBNGHWTIfJxrtunTLbJSV8kQ2lG+Qnae+/htgPhntL+pI3mrmUFkaX7f//N/fsxBkMbW6HuinVdbKVa95ZIKZRluErk7k9AmsWvolM2Naj9O88miohuDtzEf73DRNGK+ms5AWhnmh7+V3gyjxQlAhEtM5wJElBBlIvPKxcKUWwa0fjLDwacTCIEZiuklhx6nn8BpKYZ08C10dx78BdRnGpLndSRWJ3XnYq6DtflaVHN2n/dXE+jpYY+SC9G40HQRpcveR/lT35SrOec8VJFmt/Fb2LNe8679Yh67yrog/4MYBH49RMi35UjnyjVsuEzYhZvTHXOIlofTmuto/oPBsD2la//udVODeP+Tk1tw==","nonce":"KUIl6mDJNGBi","original_type":"transaction"},"id":"b74f04ce-52a0-59f3-a10a-c866448a773b","create_time":"2024-09-14T15:49:21+08:00"}`
	str := `{"appid": "wx9558d788c3884def", "mchid": "1680996035", "payer": {"openid": "oXCsa7XuDRXC1CA54oKBV5KDu32E"}, "amount": {"total": 4800, "currency": "CNY", "payer_total": 4800, "payer_currency": "CNY"}, "attach": "自定义数据说明", "bank_type": "OTHERS", "trade_type": "JSAPI", "trade_state": "SUCCESS", "out_trade_no": "1726376928081991", "success_time": "2024-09-15T13:08:59+08:00", "transaction_id": "4200002433202409154157490752", "trade_state_desc": "支付成功"}`
	val, err := utils.JsonStrToObj[payments.Transaction](str)
	if err != nil {
		panic(err)
	}
	return val
}
