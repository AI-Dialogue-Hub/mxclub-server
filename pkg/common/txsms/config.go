package txsms

type TxSmsConfig struct {
	Ak          string `yaml:"ak" validate:"required"`
	SK          string `yaml:"sk" validate:"required"`
	SmsSdkAppId string `yaml:"sms_sdk_app_id" validate:"required"`
	SignName    string `yaml:"sign_name" validate:"required"`
	TemplateId  string `yaml:"template_id" validate:"required"`
}
