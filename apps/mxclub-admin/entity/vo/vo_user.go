package vo

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"mxclub/domain/user/entity/enum"
	"time"
)

type UserLoginVO struct {
	Name     string        `json:"name"`
	Role     enum.RoleType `json:"Role"`
	JwtToken string        `json:"JwtToken"`
}

type UserVO struct {
	ID           uint             `json:"id" validate:"required" reg_err_info:"id不能为空"`
	Name         string           `json:"name"`
	WxNumber     string           `json:"wx_number"`
	WxName       string           `json:"wx_name"`
	WxIcon       string           `json:"wx_icon"`
	WxGrade      string           `json:"wx_grade"`                                     // 微信等级
	Role         enum.RoleType    `json:"role" validate:"required" reg_err_info:"不能为空"` // 用户权限
	DasherLevel  enum.DasherLevel `json:"dasher_level"`
	DisPlayName  string           `json:"disPlayName"`   // 用户权限
	MemberNumber int              `json:"member_number"` // 编号
	Phone        string           `json:"phone"`
	ActivatedAt  time.Time        `json:"activated_at"` // 最近一次上线
	CreatedAt    time.Time        `json:"created_at"`
	Bail         float64          `json:"bail"`      // 保证金
	BailTime     *time.Time       `json:"bail_time"` // 保证金缴纳时间
}

type UserTypeVO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func WrapUserTypeVOS(m maps.IMap[enum.RoleType, string]) []*UserTypeVO {
	vos := make([]*UserTypeVO, 0)
	vos = append(vos, &UserTypeVO{"all", "全部"})
	m.ForEach(func(k enum.RoleType, v string) {
		vos = append(vos, &UserTypeVO{string(k), v})
	})
	return vos
}

type AssistantOnlineVO struct {
	Id     int    `json:"id"`
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
}
