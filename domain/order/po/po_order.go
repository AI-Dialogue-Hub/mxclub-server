package po

import (
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type Order struct {
	gorm.Model
	OrderId         uint64           `gorm:"column:order_id"`
	PurchaseId      uint             `gorm:"column:purchase_id"`
	OrderName       string           `gorm:"column:order_name"`
	OrderIcon       string           `gorm:"column:icon"`
	OrderStatus     enum.OrderStatus `gorm:"column:order_status"`
	OriginalPrice   float64          `gorm:"column:original_price"`
	ProductID       uint             `gorm:"column:product_id"`
	Phone           string           `gorm:"column:phone"`
	GameRegion      string           `gorm:"column:game_region"`
	RoleId          string           `gorm:"column:role_id"`
	SpecifyExecutor bool             `gorm:"column:specify_executor"`
	ExecutorID      uint             `gorm:"column:executor_id"`
	Executor2Id     uint             `gorm:"column:executor2_id"`
	Executor3Id     uint             `gorm:"column:executor3_id"`
	Executor2Name   string           `gorm:"column:executor2_name"`
	Executor3Name   string           `gorm:"column:executor3_name"`
	Notes           string           `gorm:"column:notes"`
	DiscountPrice   float64          `gorm:"column:discount_price"`
	FinalPrice      float64          `gorm:"column:final_price"`
	ExecutorPrice   float64          `gorm:"column:executor_price"`
	Executor2Price  float64          `gorm:"column:executor_2_price"`
	Executor3Price  float64          `gorm:"column:executor_3_price"`
	PurchaseDate    *time.Time       `gorm:"column:purchase_date"`
	CompletionDate  *time.Time       `gorm:"column:completion_date"`
	DetailImages    xmysql.JSON      `gorm:"detail_images"`
}

// TableName sets the table name for the Order model.
func (Order) TableName() string {
	return "orders"
}
