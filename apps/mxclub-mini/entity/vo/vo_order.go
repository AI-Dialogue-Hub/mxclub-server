package vo

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"time"
)

type OrderVO struct {
	ID         uint   `json:"id"`
	OrderId    uint64 `json:"order_id"`
	OrderName  string `json:"order_name"`
	PurchaseId uint   `json:"purchase_id"`

	OrderStatus        enum.OrderStatus `json:"order_status"`
	OriginalPrice      float64          `json:"original_price"`
	OrderIcon          string           `json:"icon"`
	ProductID          uint             `json:"product_id"`
	Phone              string           `json:"phone"`
	GameRegion         string           `json:"game_region"`
	RoleId             string           `json:"role_id"`
	GameId             string           `json:"game_id"`
	SpecifyExecutor    bool             `json:"specify_executor"`
	ExecutorID         int              `json:"executor_id"`
	Executor2Id        int              `json:"executor2_id"`
	Executor3Id        int              `json:"executor3_id"`
	ExecutorName       string           `json:"executor_name"`
	Executor2Name      string           `json:"executor2_name"`
	Executor3Name      string           `json:"executor3_name"`
	Notes              string           `json:"notes"`
	DiscountPrice      float64          `json:"discount_price"`
	FinalPrice         float64          `json:"final_price"`
	ExecutorPrice      float64          `json:"executor_price"`
	ExecutorPriceNote  string           `json:"executor_price_note"`
	Executor2Price     float64          `json:"executor_2_price"`
	Executor2PriceNote string           `json:"executor2_price_note"`
	Executor3Price     float64          `json:"executor_3_price"`
	Executor3PriceNote string           `json:"executor3_price_note"`
	PurchaseDate       *time.Time       `json:"purchase_date"`
	CompletionDate     *time.Time       `json:"completion_date,omitempty"`
	CutRate            float64          `json:"cut_rate"`      // 抽成比例
	UserGrade          string           `json:"user_grade"`    // 老板等级
	IsEvaluation       bool             `json:"is_evaluation"` // 是否完成评价
	CanReward          bool             `json:"can_reward"`    // 是否可以打赏
}

type ProductVO struct {
	ID               uint               `json:"id"`
	Title            string             `json:"title"`
	Price            float64            `json:"price"`
	DiscountRuleID   int                `json:"discount_rule_id"`
	DiscountPrice    float64            `json:"discount_price"`
	FinalPrice       float64            `json:"final_price"`
	Description      string             `json:"description"`
	ShortDescription string             `json:"short_description"`
	Images           xmysql.StringArray `json:"images"`
	DetailImages     xmysql.StringArray `json:"detail_images"`
	Thumbnail        string             `json:"thumbnail"`
	Phone            string             `json:"phone"`
	GameId           string             `json:"gameId"`
	SalesVolume      int                `json:"sales_volume"`
	LotteryActivity  bool               `json:"lottery_activity"`
}

type WithDrawVO struct {
	HistoryWithDrawAmount float64 `json:"history_with_draw_amount"`
	WithdrawAbleAmount    float64 `json:"withdraw_able_amount"`
	WithdrawRangeMax      float64 `json:"withdraw_range_max"`
	WithdrawRangeMin      float64 `json:"withdraw_range_min"`
}

type PreferentialVO struct {
	OriginalPrice     float64 `json:"original_price"`
	DiscountedPrice   float64 `json:"discounted_price"`   // 折扣后价格
	PreferentialPrice float64 `json:"preferential_price"` // 优惠价格
	DiscountRate      float64 `json:"discount_rate"`
	DiscountInfo      string  `json:"discount_info"`
}

type WithDrawListVO struct {
	ID               uint       `json:"id"`
	PayerID          int        `json:"payer_id"`
	WithdrawalAmount float64    `json:"withdrawal_amount"`
	WithdrawalStatus string     `json:"withdrawal_status"`
	ApplicationTime  *time.Time `json:"application_time"`
	PaymentTime      *time.Time `json:"payment_time"`
	WithdrawalMethod string     `json:"withdrawal_method"`
	CreatedAt        time.Time  `json:"created_at"`
}
