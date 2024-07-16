package vo

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type OrderVO struct {
	ID              uint             `json:"id"`
	OrderId         uint64           `json:"order_id"`
	OrderName       string           `json:"order_name"`
	OrderStatus     enum.OrderStatus `json:"order_status"`
	OriginalPrice   float64          `json:"original_price"`
	OrderIcon       string           `json:"icon"`
	ProductID       uint             `json:"product_id"`
	Phone           string           `json:"phone"`
	GameRegion      string           `json:"game_region"`
	RoleId          string           `json:"role_id"`
	SpecifyExecutor bool             `json:"specify_executor"`
	ExecutorID      uint             `json:"executor_id"`
	Executor2Id     uint             `json:"executor2_id"`
	Executor3Id     uint             `json:"executor3_id"`
	Executor2Name   string           `json:"executor2_name"`
	Executor3Name   string           `json:"executor3_name"`
	Notes           string           `json:"notes"`
	DiscountPrice   float64          `json:"discount_price"`
	FinalPrice      float64          `json:"final_price"`
	ExecutorPrice   float64          `json:"executor_price"`
	Executor2Price  float64          `json:"executor_2_price"`
	Executor3Price  float64          `json:"executor_3_price"`
	PurchaseDate    *time.Time       `json:"purchase_date"`
	CompletionDate  *time.Time       `json:"completion_date,omitempty"`
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

type WithDrawVO struct {
	HistoryWithDrawAmount float64 `json:"history_with_draw_amount"`
	WithdrawAbleAmount    float64 `json:"withdraw_able_amount"`
	WithdrawRangeMax      float64 `json:"withdraw_range_max"`
	WithdrawRangeMin      float64 `json:"withdraw_range_min"`
}

type PreferentialVO struct {
	OriginalPrice     float64 `json:"original_price"`
	DiscountedPrice   float64 `json:"discounted_price"`
	PreferentialPrice float64 `json:"preferential_price"`
	DiscountRate      float64 `json:"discount_rate"`
	DiscountInfo      string  `json:"discount_info"`
}
