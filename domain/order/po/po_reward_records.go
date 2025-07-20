package po

import (
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"time"
)

// RewardRecord 表示reward_records表的结构体
type RewardRecord struct {
	ID           uint             `gorm:"primaryKey;autoIncrement:true" json:"id"`           // 自增主键
	PurchaserID  uint             `gorm:"not null" json:"purchaser_id"`                      // 购买人的ID
	OrderID      string           `gorm:"uniqueIndex:idx_order_id;not null" json:"order_id"` // 打赏订单Id
	DasherID     int              `gorm:"not null" json:"dasher_id"`                         // 打手的ID
	DasherNumber uint             `gorm:"not null" json:"dasher_number"`                     // 打手db id
	DasherName   string           `gorm:"size:50;not null" json:"dasher_name"`               // 打手name
	Remarks      string           `gorm:"type:text" json:"remarks"`                          // 备注信息
	RewardAmount float64          `gorm:"type:decimal(10,2);not null" json:"reward_amount"`  // 打赏金额
	RewardTime   time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"reward_time"`      // 打赏时间
	Status       enum.OrderStatus `gorm:"not null" json:"status"`
	OutTradeNo   string           `gorm:"size:50;not null" json:"out_trade_no"` // 微信支付唯一标识
	CreatedAt    time.Time        `json:"created_at"`                           // 创建时间
	UpdatedAt    time.Time        `json:"updated_at"`                           // 更新时间
	DeletedAt    gorm.DeletedAt   `json:"deleted_at,omitempty"`                 // 删除时间
}

func (RewardRecord) TableName() string {
	return "reward_records"
}
