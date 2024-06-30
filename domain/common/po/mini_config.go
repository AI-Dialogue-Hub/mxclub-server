// 小程序里面的配置信息
package po

import (
	"github.com/jinzhu/gorm"
)

type MiniConfig struct {
	gorm.Model
	ConfigName string `gorm:"column:config_name;size:50;not null"`
	Content    []byte `gorm:"column:content;type:json"`
	Role       uint64 `gorm:"size:20;not null"`
}

func (u *MiniConfig) TableName() string {
	return "mx_mini_config"
}
