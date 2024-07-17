package req

type CaptchaReq struct {
	CaptchaId string `json:"captcha_id" form:"captcha_id" validate:"required"`
	Answer    string `json:"answer" form:"answer" validate:"required"`
}

type AssistantReq struct {
	MemberNumber int64  `json:"member_number" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Name         string `json:"name"`
}

type UserInfoReq struct {
	AvatarUrlBase64 string `json:"avatarUrl"` // base64
}
