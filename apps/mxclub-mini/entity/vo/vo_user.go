package vo

import (
	"mxclub/domain/user/entity/enum"
)

type User struct {
	WxNumber string        `json:"wx_number"`
	WxName   string        `json:"wx_name"`
	WxIcon   string        `json:"wx_icon,omitempty"`
	WxGrade  string        `json:"wx_grade"`
	Role     enum.RoleType `json:"role"`
}

type CaptchaVO struct {
	CaptchaId string `json:"captcha_id"`
	B64s      string `json:"b64_s"`
	Answer    string `json:"answer"`
}

type CaptchaVerifyVO struct {
	CaptchaId string `json:"captcha_id"`
	Answer    string `json:"answer"`
	Result    bool   `json:"result"`
}

type AssistantOnlineVO struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
