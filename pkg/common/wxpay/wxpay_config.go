package wxpay

type WxPayConfig struct {
	MchID                      string `yaml:"mch_id"`                        // 商户号
	MchCertificateSerialNumber string `yaml:"mch_certificate_serial_number"` // 商户证书序列号
	MchAPIv3Key                string `yaml:"mch_ap_iv_3_key"`               // 商户APIv3密钥
	PrivateKeyPath             string `yaml:"private_key_path"`
	PublicKeyPath              string `yaml:"public_key_path"` // 新的商户使用这个
	WechatpayPublicKeyID       string `yaml:"wechatpay_public_key_id"`
	AppId                      string `yaml:"app_id"`
	PayName                    string `yaml:"pay_name"`              // 是哪个小程序在使用
	CallBackURL                string `yaml:"call_back_url"`         // 支付回调接口
	RefundsCallBackURL         string `yaml:"refunds_call_back_url"` // 退款回调接口
}

func (c *WxPayConfig) IsBaoZaoClub() bool {
	return c != nil && c.PayName == "baozao-club"
}
