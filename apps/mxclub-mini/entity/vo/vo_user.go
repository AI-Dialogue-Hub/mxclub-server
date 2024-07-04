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
