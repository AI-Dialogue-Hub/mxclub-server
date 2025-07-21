package po

import (
	"gorm.io/gorm"
	"time"
)

type DeactivateDasher struct {
	ID                    uint           `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	DasherID              int            `gorm:"column:dasher_id;default:-1;comment:打手编号"`
	DasherName            string         `gorm:"column:dasher_name;type:varchar(100);not null;comment:打手姓名"`
	HistoryWithdrawAmount float64        `gorm:"column:history_with_draw_amount;type:decimal(10,2);comment:历史提现的钱"`
	WithdrawAbleAmount    float64        `gorm:"column:withdraw_able_amount;type:decimal(10,2);comment:还能提现的钱"`
	OrderSnapshot         string         `gorm:"column:order_snapshot;type:text;comment:gzip压缩的订单快照信息"`
	CreatedAt             time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt             *time.Time     `gorm:"column:updated_at;type:timestamp;default:NULL;comment:更新时间"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index;default:NULL;comment:删除时间(软删除)"`
}

// TableName 设置表名
func (DeactivateDasher) TableName() string {
	return "deactivate_dasher"
}
