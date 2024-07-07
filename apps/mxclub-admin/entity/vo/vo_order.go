package vo

import (
	"mxclub/domain/product/entity/enum"
	"time"
)

type OrderVO struct {
	ID              uint             `json:"id,omitempty"`
	PurchaseId      uint             `json:"purchase_id,omitempty"`
	OrderName       string           `json:"order_name,omitempty"`
	OrderIcon       string           `json:"order_icon,omitempty"`
	OrderStatus     enum.OrderStatus `json:"order_status,omitempty"`
	OrderStatusStr  string           `json:"order_status_str,omitempty"`
	OriginalPrice   float64          `json:"original_price,omitempty"`
	ProductID       uint             `json:"product_id,omitempty"`
	GameRegion      string           `json:"game_region,omitempty"`
	SpecifyExecutor bool             `json:"specify_executor,omitempty"`
	ExecutorID      uint             `json:"executor_id,omitempty"`
	Notes           string           `json:"notes,omitempty"`
	DiscountPrice   float64          `json:"discount_price,omitempty"`
	FinalPrice      float64          `json:"final_price,omitempty"`
	PurchaseDate    *time.Time       `json:"purchase_date,omitempty"`
	CompletionDate  *time.Time       `json:"completion_date,omitempty"`
}
