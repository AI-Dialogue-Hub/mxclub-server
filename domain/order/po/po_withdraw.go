package po

import (
	"gorm.io/gorm"
	"time"
)

type WithdrawalRecord struct {
	ID               uint    `gorm:"primaryKey"`
	DasherID         int     `gorm:"not null"`
	DasherUserId     uint    `gorm:"not null;column:dasher_user_id"`
	DasherName       string  `gorm:"not null"`
	PayerID          int     `gorm:"default:null"`
	WithdrawalAmount float64 `gorm:"not null"`
	WithdrawalStatus string  `gorm:"type:enum('initiated', 'completed', 'reject'); not null"`
	// 提现申请时间
	ApplicationTime  *time.Time `gorm:"not null"`
	PaymentTime      *time.Time `gorm:"default:null"`
	WithdrawalMethod string     `gorm:"size:100"`
	CreatedAt        time.Time  `gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt
}

// TableName sets the table name for the WithdrawalRecord model.
func (WithdrawalRecord) TableName() string {
	return "withdrawal_records"
}
