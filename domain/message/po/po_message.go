package po

import (
	"gorm.io/gorm"
	"mxclub/domain/message/entity/enum"
	"time"
)

type Message struct {
	ID int `gorm:"primaryKey"`

	MessageType   enum.MessageType   `gorm:"type:int;not null"`
	Title         string             `gorm:"type:varchar(128);not null"`
	Content       string             `gorm:"type:varchar(512)"`
	MessageFrom   int                `gorm:"not null"`
	MessageTo     uint               `gorm:"not null"`
	MessageStatus enum.MessageStatus `gorm:"not null"` // 0 未读 1已读
	Ext           string             `gorm:"type:varchar(50)"`

	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Message) TableName() string {
	return "messages"
}
