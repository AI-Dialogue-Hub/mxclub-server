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
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	Gender    int    `json:"gender"` // 使用 int 类型来表示 Number 类型
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
}
