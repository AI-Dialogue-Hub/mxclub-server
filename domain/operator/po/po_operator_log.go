package po

import (
	"mxclub/domain/operator/entity/enum"
	"time"
)

type OperatorLogPO struct {
	ID uint `gorm:"primaryKey"`

	Type     enum.OperatorEnum `gorm:"column:type"`
	Remarks  string
	UserId   uint   `gorm:"column:user_id"`
	UserName string `gorm:"column:user_name"`

	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `gorm:"default:NULL"`
}

func (*OperatorLogPO) TableName() string {
	return "operator_log"
}
