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
	Name         string
	Status       string `gorm:"default:PENDING"` // PENDING PASS REJECT

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (AssistantApplication) TableName() string {
	return "assistant_application"
}
