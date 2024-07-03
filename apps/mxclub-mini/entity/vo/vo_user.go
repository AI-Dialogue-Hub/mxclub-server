package vo

import (
	"mxclub/domain/user/entity/enum"
)

type User struct {
	WxNumber string        `json:"wx_number"`
	WxName   string        `json:"wx_name"`
	Role     enum.RoleType `json:"role"`
}
