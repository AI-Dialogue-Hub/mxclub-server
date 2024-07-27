package vo

import (
	"mxclub/domain/user/entity/enum"
)

type UserVO struct {
	WxNumber        string        `json:"wx_number"`
	WxName          string        `json:"wx_name"`
	WxIcon          string        `json:"wx_icon,omitempty"`
	WxGrade         string        `json:"wx_grade"`
	Role            enum.RoleType `json:"role"`
	MemberNumber    int           `json:"member_number"`
	CurrentPoints   float64       `json:"currentPoints"`
	NextLevelPoints int           `json:"nextLevelPoints"`
}

func (userVO *UserVO) SetCurrentPoints(currentPoints float64) {
	userVO.CurrentPoints = currentPoints
	for _, amount := range gradeRule {
		if currentPoints < amount {
			userVO.NextLevelPoints = int(amount)
			break
		}
	}
}

var gradeRule = []float64{500, 2000, 5000, 10000, 20000, 50000}

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
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
