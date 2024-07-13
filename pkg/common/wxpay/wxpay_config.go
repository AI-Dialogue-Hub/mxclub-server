package wxpay

type WxPayConfig struct {
	MchID                      string `yaml:"mch_id"`                        // 商户号
	MchCertificateSerialNumber string `yaml:"mch_certificate_serial_number"` // 商户证书序列号
	MchAPIv3Key                string `yaml:"mch_ap_iv_3_key"`               // 商户APIv3密钥
	PrivateKeyPath             string `yaml:"private_key_path"`
	AppId                      string `yaml:"app_id"`
}
