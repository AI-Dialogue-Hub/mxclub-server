package req

type CaptchaReq struct {
	CaptchaId string `json:"captcha_id" form:"captcha_id" validate:"required"`
	Answer    string `json:"answer" form:"answer" validate:"required"`
}
