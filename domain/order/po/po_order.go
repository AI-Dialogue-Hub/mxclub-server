package po

import (
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type Order struct {
	gorm.Model
	OrderId            uint64           `gorm:"column:order_id"`
	PurchaseId         uint             `gorm:"column:purchase_id"`
	OrderName          string           `gorm:"column:order_name"`
	OrderIcon          string           `gorm:"column:icon"`
	OrderStatus        enum.OrderStatus `gorm:"column:order_status"`
	OriginalPrice      float64          `gorm:"column:original_price"`
	ProductID          uint             `gorm:"column:product_id"`
	Phone              string           `gorm:"column:phone"`
	GameRegion         string           `gorm:"column:game_region"`
	RoleId             string           `gorm:"column:role_id"`
	SpecifyExecutor    bool             `gorm:"column:specify_executor"`
	ExecutorID         uint             `gorm:"column:executor_id"`
	ExecutorName       string           `gorm:"column:executor_name"`
	Executor2Id        uint             `gorm:"column:executor2_id"`
	Executor3Id        uint             `gorm:"column:executor3_id"`
	Executor2Name      string           `gorm:"column:executor2_name"`
	Executor3Name      string           `gorm:"column:executor3_name"`
	Notes              string           `gorm:"column:notes"`
	DiscountPrice      float64          `gorm:"column:discount_price"`
	FinalPrice         float64          `gorm:"column:final_price"`
	ExecutorPrice      float64          `gorm:"column:executor_price"`
	ExecutorPriceNote  string           `gorm:"column:executor_price_note"`
	Executor2Price     float64          `gorm:"column:executor2_price"`
	Executor2PriceNote string           `gorm:"column:executor2_price_note"`
	Executor3Price     float64          `gorm:"column:executor3_price"`
	Executor3PriceNote string           `gorm:"column:executor3_price_note"`
	PurchaseDate       *time.Time       `gorm:"column:purchase_date"`
	CompletionDate     *time.Time       `gorm:"column:completion_date"`
	DetailImages       xmysql.JSON      `gorm:"detail_images"`        // 订单结束后上传的图片
	CutRate            float64          `gorm:"column:cut_rate"`      // 抽成比例
	IsEvaluation       bool             `gorm:"column:is_evaluation"` // 是否完成评价
}

// TableName sets the table name for the Order model.
func (Order) TableName() string {
	return "orders"
}
