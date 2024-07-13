package req

import (
	"mxclub/domain/user/entity/enum"
)

type UserReq struct {
	ID           uint          `json:"id" validate:"required" reg_err_info:"id不能为空"`
	Name         string        `json:"name"`
	WxName       string        `json:"wx_name"`
	WxIcon       string        `json:"wx_icon"`
	WxGrade      string        `json:"wx_grade"`                                     // 微信等级
	Role         enum.RoleType `json:"role" validate:"required" reg_err_info:"不能为空"` // 用户权限
	MemberNumber string        `json:"member_number"`                                // 编号
}