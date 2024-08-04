package po

import "time"

type OrderEvaluation struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	OrdersID   uint   `gorm:"not null"`  // 订单的 orders_id
	OrderID    uint   `gorm:"not null"`  // 订单的 order_id
	ExecutorID uint   `gorm:"not null"`  // 打手 id
	Rating     int    `gorm:"not null"`  // 评价的评分（假设是从 1 到 5）
	Comments   string `gorm:"type:text"` // 评价的评论

	CreatedAt time.Time  `gorm:"autoCreateTime"` // 创建时间
	UpdatedAt time.Time  `gorm:"autoUpdateTime"` // 更新时间
	DeletedAt *time.Time `gorm:"default:NULL"`
}

func (*OrderEvaluation) TableName() string {
	return "order_evaluation"
}