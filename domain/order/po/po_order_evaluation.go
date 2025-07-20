package po

import (
	"gorm.io/gorm"
	"time"
)

type OrderEvaluation struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	OrdersID   uint   `gorm:"not null"`  // 订单的 orders_id  => db id
	OrderID    uint64 `gorm:"not null"`  // 订单的 order_id
	ExecutorID int    `gorm:"not null"`  // 打手 id
	Rating     int    `gorm:"not null"`  // 评价的评分（假设是从 1 到 5）
	Comments   string `gorm:"type:text"` // 评价的评论

	CreatedAt time.Time `gorm:"autoCreateTime"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // 更新时间
	DeletedAt gorm.DeletedAt
}

func (*OrderEvaluation) TableName() string {
	return "order_evaluation"
}
