package po

import (
	"database/sql"
	"gorm.io/gorm"
	"mxclub/domain/user/entity/enum"
)

type User struct {
	gorm.Model
	Name         string            `gorm:"name"`
	WxNumber     string            `gorm:"wx_number"`
	WxName       string            `gorm:"wx_name"`
	WxOpenId     string            `gorm:"wx_open_id"`
	WxIcon       string            `gorm:"wx_icon"`
	WxGrade      string            `gorm:"wx_grade"`
	WxUserInfo   string            `gorm:"wx_user_info"`
	Role         enum.RoleType     `gorm:"role"`
	Phone        string            `gorm:"role"`
	MemberNumber uint              `gorm:"member_number"` // 打手编号
	MemberStatus enum.MemberStatus `gorm:"member_status"` // 编号
	ActivatedAt  sql.NullTime      `gorm:"activated_at"`  // 最近一次上线
}

func (u User) TableName() string {
	return "mx_user"
}
