package po

import (
	"mxclub/domain/order/entity/enum"
	"time"
)

type Deduction struct {
	ID              uint `gorm:"primaryKey"`
	UserID          uint
	OrderNo         uint
	ConfirmPersonId uint
	Amount          float64
	Reason          string
	Status          enum.DeductStatus `gorm:"default:'PENDING'"`
	CreatedAt       time.Time         `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time         `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt       *time.Time        `gorm:"default:NULL"`
}

func (*Deduction) TableName() string {
	return "order_deduct"
}
