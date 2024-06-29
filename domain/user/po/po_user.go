package po

import (
	"database/sql"
	"gorm.io/gorm"
	"mxclub/domain/user/entity/enum"
)

type User struct {
	gorm.Model
	Name         string         `gorm:"name"`
	WxNumber     string         `gorm:"wx_number"`
	Role         enum.RoleType  `gorm:"role"`
	MemberNumber sql.NullString `gorm:"member_number"` // 编号
	ActivatedAt  sql.NullTime   `gorm:"activated_at"`  // 最近一次上线
}

func (u User) TableName() string {
	return "mx_user"
}
