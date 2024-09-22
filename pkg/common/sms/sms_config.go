package sms

type Config struct {
	Ak        string `yaml:"ak" validate:"required"`
	Sk        string `yaml:"sk" validate:"required"`
	SmsCode   string `yaml:"sms_code" validate:"required"`
	SignName  string `yaml:"sign_name" validate:"required"`
	TestPhone string `yaml:"test_phone" validate:"required"`
}

type DispatchReq struct {
	Consignee string `json:"consignee"`
	Role      string `json:"role"`
}

func NewDispatchReq(consignee, role string) *DispatchReq {
	return &DispatchReq{
		Consignee: consignee,
		Role:      role,
	}
}
