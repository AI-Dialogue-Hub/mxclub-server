package vo

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type OrderVO struct {
	OrderName       string           `json:"order_name"`
	OrderStatus     enum.OrderStatus `json:"order_status"`
	OriginalPrice   float64          `json:"original_price"`
	OrderIcon       string           `json:"icon"`
	ProductID       uint             `json:"product_id"`
	GameRegion      string           `json:"game_region"`
	SpecifyExecutor bool             `json:"specify_executor"`
	ExecutorID      uint             `json:"executor_id"`
	Notes           string           `json:"notes"`
	DiscountPrice   float64          `json:"discount_price"`
	FinalPrice      float64          `json:"final_price"`
	PurchaseDate    *time.Time       `json:"purchase_date"`
	CompletionDate  *time.Time       `json:"completion_date"`
}

type ProductVO struct {
	ID               uint        `json:"id"`
	Title            string      `json:"title"`
	Price            float64     `json:"price"`
	DiscountRuleID   int         `json:"discount_rule_id"`
	DiscountPrice    float64     `json:"discount_price"`
	FinalPrice       float64     `json:"final_price"`
	Description      string      `json:"description"`
	ShortDescription string      `json:"short_description"`
	Images           xmysql.JSON `json:"images"`
	DetailImages     xmysql.JSON `json:"detail_images"`
}
