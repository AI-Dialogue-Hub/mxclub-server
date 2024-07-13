package wxpay

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	Client          *core.Client
	paymentSVC      *app.AppApiService
	wxPayConfig     *WxPayConfig
	jsapiApiService *jsapi.JsapiApiService
)

func NewWxPayClient(config *WxPayConfig) *core.Client {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(config.PrivateKeyPath)
	if err != nil {
		xlog.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(config.MchID, config.MchCertificateSerialNumber, mchPrivateKey, config.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		xlog.Fatalf("new wechat pay client err:%s", err)
	}
	Client = client
	jsapiApiService = &jsapi.JsapiApiService{Client: client}
	wxPayConfig = config
	return client
}
