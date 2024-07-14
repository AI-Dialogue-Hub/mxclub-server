package vo

import (
	"mxclub/domain/user/entity/enum"
	"time"
)

type UserLoginVO struct {
	Name     string        `json:"name"`
	Role     enum.RoleType `json:"Role"`
	JwtToken string        `json:"JwtToken"`
}

type UserVO struct {
	ID           uint          `json:"id" validate:"required" reg_err_info:"id不能为空"`
	Name         string        `json:"name"`
	WxNumber     string        `json:"wx_number"`
	WxName       string        `json:"wx_name"`
	WxIcon       string        `json:"wx_icon"`
	WxGrade      string        `json:"wx_grade"`                                     // 微信等级
	Role         enum.RoleType `json:"role" validate:"required" reg_err_info:"不能为空"` // 用户权限
	DisPlayName  string        `json:"disPlayName"`                                  // 用户权限
	MemberNumber int           `json:"member_number"`                                // 编号
	ActivatedAt  time.Time     `json:"activated_at"`                                 // 最近一次上线
}
