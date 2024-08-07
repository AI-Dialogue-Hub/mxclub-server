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

func TestRefunds(t *testing.T) {
	info := getInfo(t)
	requestCtx := &fasthttp.RequestCtx{Request: fasthttp.Request{}, Response: fasthttp.Response{}}
	ctx := jetContext.NewContext(requestCtx, xlog.NewWith("text"))
	err := Refunds(ctx, info, GenerateOutRefundNo(), "商品已售完")
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
