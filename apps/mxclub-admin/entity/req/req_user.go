package req

import (
	"mxclub/domain/user/entity/enum"
	"mxclub/pkg/api"
)

type UserReq struct {
	ID           uint          `json:"id" validate:"required" reg_err_info:"id不能为空"`
	Name         string        `json:"name"`
	WxName       string        `json:"wx_name"`
	WxIcon       string        `json:"wx_icon"`
	WxGrade      string        `json:"wx_grade"`                                     // 微信等级
	Role         enum.RoleType `json:"role" validate:"required" reg_err_info:"不能为空"` // 用户权限
	MemberNumber any           `json:"member_number"`                                // 编号
	Phone        string        `json:"phone"`
}

type UserListReq struct {
	*api.PageParams
	UserType     string `json:"role"` // 用户类型
	MemberNumber int    `json:"memberNumber"`
}

type UserRemoveReq struct {
	DasherIds []int `json:"dasher_ids" validate:"required" reg_err_info:"不能为空"`
	OrderId   uint  `json:"order_id" validate:"required" reg_err_info:"不能为空"`
}
