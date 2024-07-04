package vo

import (
	"mxclub/domain/order/entity/enum"
	"time"
)

type OrderVO struct {
	OrderName       string           `json:"order_name"`
	OrderStatus     enum.OrderStatus `json:"order_status"`
	OriginalPrice   float64          `json:"original_price"`
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
