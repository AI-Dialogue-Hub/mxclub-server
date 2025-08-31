package po

import (
	"database/sql"
	"gorm.io/gorm"
	"mxclub/domain/user/entity/enum"
)

type User struct {
	gorm.Model
	Name         string            `gorm:"name"` // 打手名字
	Password     string            `gorm:"password"`
	WxNumber     string            `gorm:"wx_number"`
	WxName       string            `gorm:"wx_name"`
	WxOpenId     string            `gorm:"wx_open_id"`
	WxIcon       string            `gorm:"wx_icon"`
	WxGrade      string            `gorm:"wx_grade"`
	WxUserInfo   string            `gorm:"wx_user_info"`
	Role         enum.RoleType     `gorm:"role"`
	Phone        string            `gorm:"role"`
	GameId       string            `gorm:"game_id"`                 // 游戏Id
	MemberNumber int               `gorm:"member_number"`           // 打手编号
	DasherLevel  enum.DasherLevel  `gorm:"dasher_level;default:-1"` // 打手等级
	MemberStatus enum.MemberStatus `gorm:"member_status"`           // 编号
	ActivatedAt  sql.NullTime      `gorm:"activated_at"`            // 最近一次上线
}

func (u User) TableName() string {
	return "mx_user"
}
