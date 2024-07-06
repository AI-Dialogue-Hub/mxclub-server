package po

import (
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"time"
)

type Order struct {
	gorm.Model
	PurchaseId      uint             `gorm:"column:purchase_id"`
	OrderName       string           `gorm:"column:order_name"`
	OrderIcon       string           `gorm:"column:icon"`
	OrderStatus     enum.OrderStatus `gorm:"column:order_status"`
	OriginalPrice   float64          `gorm:"column:original_price"`
	ProductID       uint             `gorm:"column:product_id"`
	GameRegion      string           `gorm:"column:game_region"`
	SpecifyExecutor bool             `gorm:"column:specify_executor"`
	ExecutorID      uint             `gorm:"column:executor_id"`
	Notes           string           `gorm:"column:notes"`
	DiscountPrice   float64          `gorm:"column:discount_price"`
	FinalPrice      float64          `gorm:"column:final_price"`
	PurchaseDate    *time.Time       `gorm:"column:purchase_date"`
	CompletionDate  *time.Time       `gorm:"column:completion_date"`
}

// TableName sets the table name for the Order model.
func (Order) TableName() string {
	return "orders"
}
