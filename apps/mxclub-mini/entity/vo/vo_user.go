package vo

import (
	"mxclub/domain/user/entity/enum"
)

type User struct {
	Name     string        `gorm:"name"`
	WxNumber string        `gorm:"wx_number"`
	Role     enum.RoleType `gorm:"role"`
}
