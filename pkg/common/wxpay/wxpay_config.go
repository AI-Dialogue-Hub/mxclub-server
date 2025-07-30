package wxpay

type WxPayConfig struct {
	MchID                      string `yaml:"mch_id"`                        // 商户号
	MchCertificateSerialNumber string `yaml:"mch_certificate_serial_number"` // 商户证书序列号
	MchAPIv3Key                string `yaml:"mch_ap_iv_3_key"`               // 商户APIv3密钥
	PrivateKeyPath             string `yaml:"private_key_path"`
	PublicKeyPath              string `yaml:"public_key_path"`         // 新的商户使用这个
	WechatpayPublicKeyID       string `yaml:"wechatpay_public_key_id"` // 公钥id
	AppId                      string `yaml:"app_id"`
	PayName                    string `yaml:"pay_name"`              // 是哪个小程序在使用
	NewPay                     bool   `yaml:"new_pay"`               // 是否使用新的支付方式
	ClubAvatar                 string `yaml:"club_avatar"`           // 小程序默认头像
	CallBackURL                string `yaml:"call_back_url"`         // 支付回调接口
	RefundsCallBackURL         string `yaml:"refunds_call_back_url"` // 退款回调接口
	RewardCallBackURL          string `yaml:"reward_call_back_url"`  // 打赏支付支付回调接口
}

func (c *WxPayConfig) IsBaoZaoClub() bool {
	return c != nil && c.PayName == "baozao-club"
}

// IsNewPay 是否是新的支付方式，除明星电竞外都是
func (c *WxPayConfig) IsNewPay() bool {
	return c != nil && c.NewPay
}
