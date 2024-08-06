package wxpay

import (
	"context"
	"crypto/rsa"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	Client            *core.Client
	wxPayConfig       *WxPayConfig
	jsapiApiService   *jsapi.JsapiApiService
	mchPrivateKey     *rsa.PrivateKey
	notifyHandler     *Handler
	refundsApiService *refunddomestic.RefundsApiService
)

func InitWxPay(config *WxPayConfig) {
	NewWxPayClient(config)
	NewWxPayCertHandler(config)
	NewRefundsApiService()
}

func NewWxPayClient(config *WxPayConfig) *core.Client {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	var err error
	mchPrivateKey, err = utils.LoadPrivateKeyWithPath(config.PrivateKeyPath)
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

// NewWxPayCertHandler 处理回调参数
func NewWxPayCertHandler(config *WxPayConfig) *Handler {
	var (
		ctx = context.Background()
		err error
	)
	if mchPrivateKey == nil {
		mchPrivateKey, err = utils.LoadPrivateKeyWithPath(config.PrivateKeyPath)
		if err != nil {
			xlog.Fatal("load merchant private key error", err)
		}
	}
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(
		ctx, mchPrivateKey, config.MchCertificateSerialNumber, config.MchID, config.MchAPIv3Key,
	)
	if err != nil {
		xlog.Fatal("RegisterDownloaderWithPrivateKey error", err)
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(config.MchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, err := NewRSANotifyHandler(config.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		xlog.Fatal("NewRSANotifyHandler error", err)
	}
	return handler
}

func NewRefundsApiService() *refunddomestic.RefundsApiService {
	refundsApiService = &refunddomestic.RefundsApiService{Client: Client}
	return refundsApiService
}
