package vo

import (
	"mxclub/domain/order/entity/enum"
	"time"
)

type RewardVO struct {
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
	CreatedAt    time.Time        `json:"created_at"` // 创建时间
}
