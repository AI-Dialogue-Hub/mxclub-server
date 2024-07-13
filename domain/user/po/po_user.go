package po

import (
	"database/sql"
	"gorm.io/gorm"
	"mxclub/domain/user/entity/enum"
)

type User struct {
	gorm.Model
	Name         string        `gorm:"name"`
	WxNumber     string        `gorm:"wx_number"`
	WxName       string        `gorm:"wx_name"`
	WxOpenId     string        `gorm:"wx_open_id"`
	WxIcon       string        `gorm:"wx_icon"`
	WxGrade      string        `gorm:"wx_grade"`
	Role         enum.RoleType `gorm:"role"`
	Phone        string        `gorm:"role"`
	MemberNumber int           `gorm:"member_number"` // 编号
	ActivatedAt  sql.NullTime  `gorm:"activated_at"`  // 最近一次上线
}

func (u User) TableName() string {
	return "mx_user"
}
