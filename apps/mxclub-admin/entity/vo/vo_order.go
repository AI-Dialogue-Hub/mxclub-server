package vo

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type OrderVO struct {
	ID                 uint             `json:"id,omitempty"`
	OrderId            uint64           `json:"order_id"`
	PurchaseId         uint             `json:"purchase_id"`
	OrderName          string           `json:"order_name"`
	OrderIcon          string           `json:"icon"`
	OrderStatus        enum.OrderStatus `json:"order_status"`
	OrderStatusStr     string           `json:"order_status_str"`
	OriginalPrice      float64          `json:"original_price"`
	ProductID          uint             `json:"product_id"`
	Phone              string           `json:"phone"`
	GameRegion         string           `json:"game_region"`
	RoleId             string           `json:"role_id"`
	SpecifyExecutor    bool             `json:"specify_executor"`
	ExecutorID         uint             `json:"executor_id"`
	ExecutorName       string           `json:"executor_name"`
	Executor2Id        uint             `json:"executor2_id"`
	Executor3Id        uint             `json:"executor3_id"`
	Executor2Name      string           `json:"executor2_name"`
	Executor3Name      string           `json:"executor3_name"`
	Notes              string           `json:"notes"`
	DiscountPrice      float64          `json:"discount_price"`
	FinalPrice         float64          `json:"final_price"`
	ExecutorPrice      float64          `json:"executor_price"`
	ExecutorPriceNote  string           `json:"executor_price_note"`
	Executor2Price     float64          `json:"executor2_price"`
	Executor2PriceNote string           `json:"executor2_price_note"`
	Executor3Price     float64          `json:"executor3_price"`
	Executor3PriceNote string           `json:"executor3_price_note"`
	PurchaseDate       *time.Time       `json:"purchase_date"`
	CompletionDate     *time.Time       `json:"completion_date"`
	DetailImages       xmysql.JSON      `json:"detail_images"` // 订单结束后上传的图片
	CutRate            float64          `json:"cut_rate"`      // 抽成比例
}

type WithdrawVO struct {
	ID               uint       `json:"id"`
	DasherID         uint       `json:"dasher_id"`
	DasherName       string     `json:"dasher_name"`
	PayerID          int        `json:"payer_id"`
	WithdrawalAmount float64    `json:"withdrawal_amount"`
	WithdrawalStatus string     `json:"withdrawal_status"`
	ApplicationTime  *time.Time `json:"application_time"`
	PaymentTime      *time.Time `json:"payment_time"`
	WithdrawalMethod string     `json:"withdrawal_method"`
	CreatedAt        time.Time  `json:"created_at"`
}
