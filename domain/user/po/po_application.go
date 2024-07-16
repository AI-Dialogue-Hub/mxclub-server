package po

import (
	"gorm.io/gorm"
	"time"
)

type AssistantApplication struct {
	ID           uint `gorm:"primarykey"`
	UserID       uint
	Phone        string
	MemberNumber int64
	Status       string `gorm:"default:PENDING"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
