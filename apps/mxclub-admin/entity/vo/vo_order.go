package vo

import (
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type OrderVO struct {
	ID                 uint                  `json:"id,omitempty"`
	OrderId            uint64                `json:"order_id"`
	PurchaseId         uint                  `json:"purchase_id"`
	OrderName          string                `json:"order_name"`
	OrderIcon          string                `json:"icon"`
	OrderStatus        enum.OrderStatus      `json:"order_status"`
	OrderStatusStr     string                `json:"order_status_str"`
	OriginalPrice      float64               `json:"original_price"`
	ProductID          uint                  `json:"product_id"`
	Phone              string                `json:"phone"`
	GameRegion         string                `json:"game_region"`
	RoleId             string                `json:"role_id"`
	SpecifyExecutor    bool                  `json:"specify_executor"`
	ExecutorID         int                   `json:"executor_id"`
	ExecutorName       string                `json:"executor_name"`
	Executor2Id        int                   `json:"executor2_id"`
	Executor3Id        int                   `json:"executor3_id"`
	Executor2Name      string                `json:"executor2_name"`
	Executor3Name      string                `json:"executor3_name"`
	Notes              string                `json:"notes"`
	DiscountPrice      float64               `json:"discount_price"`
	FinalPrice         float64               `json:"final_price"`
	ExecutorPrice      float64               `json:"executor_price"`
	ExecutorPriceNote  string                `json:"executor_price_note"`
	Executor2Price     float64               `json:"executor2_price"`
	Executor2PriceNote string                `json:"executor2_price_note"`
	Executor3Price     float64               `json:"executor3_price"`
	Executor3PriceNote string                `json:"executor3_price_note"`
	PurchaseDate       *time.Time            `json:"purchase_date"`
	CompletionDate     *time.Time            `json:"completion_date"`
	StartImages        string                `json:"start_images"`
	DetailImages       xmysql.StringArray    `json:"detail_images"`   // 订单结束后上传的图片
	CutRate            float64               `json:"cut_rate"`        // 抽成比例
	IsEvaluation       bool                  `json:"is_evaluation"`   // 是否完成评价
	EvaluationInfo     map[int]*EvaluationVO `json:"evaluation_info"` // 评价信息 ExecutorID => evaluation
}

type EvaluationVO struct {
	ID         uint   `json:"id"`
	OrdersID   uint   `json:"orders_id"`
	OrderID    uint   `json:"order_id"`
	ExecutorID int    `json:"executor_id"`
	Rating     int    `json:"rating"`
	Comments   string `json:"comments"`
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

// WithDrawVO 历史提现的钱
type WithDrawVO struct {
	HistoryWithDrawAmount float64 `json:"history_with_draw_amount"`
	WithdrawAbleAmount    float64 `json:"withdraw_able_amount"`
	WithdrawRangeMax      float64 `json:"withdraw_range_max"`
	WithdrawRangeMin      float64 `json:"withdraw_range_min"`
}

type HistoryWithDrawVO struct {
	*WithDrawVO
	DasherID   uint   `json:"dasher_id"`
	DasherName string `json:"dasher_name"`
}

type DeactivateDasherVO struct {
	ID                    uint           `json:"id"`
	DasherID              int            `json:"dasher_id"`
	DasherName            string         `json:"dasher_name"`
	HistoryWithdrawAmount float64        `json:"history_withdraw_amount"`
	WithdrawAbleAmount    float64        `json:"withdraw_able_amount"`
	OrderSnapshot         string         `json:"order_snapshot"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             *time.Time     `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`
}
