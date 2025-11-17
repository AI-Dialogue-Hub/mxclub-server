package po

import (
	"fmt"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
	"regexp"
	"time"
)

type Order struct {
	gorm.Model
	OrderId            uint64             `gorm:"column:order_id"`
	PurchaseId         uint               `gorm:"column:purchase_id"`
	OrderName          string             `gorm:"column:order_name"`
	OrderIcon          string             `gorm:"column:icon"`
	OrderStatus        enum.OrderStatus   `gorm:"column:order_status"`
	OriginalPrice      float64            `gorm:"column:original_price"`
	ProductID          uint               `gorm:"column:product_id"`
	Phone              string             `gorm:"column:phone"`
	GameRegion         string             `gorm:"column:game_region"`
	RoleId             string             `gorm:"column:role_id"` // 这里可能包含两部分 like => Id:9913193068,角色:六哥会打呀
	SpecifyExecutor    bool               `gorm:"column:specify_executor"`
	ExecutorID         int                `gorm:"column:executor_id"`
	ExecutorName       string             `gorm:"column:executor_name"`
	Executor2Id        int                `gorm:"column:executor2_id"`
	Executor3Id        int                `gorm:"column:executor3_id"`
	Executor2Name      string             `gorm:"column:executor2_name"`
	Executor3Name      string             `gorm:"column:executor3_name"`
	Notes              string             `gorm:"column:notes"`
	DiscountPrice      float64            `gorm:"column:discount_price"`
	FinalPrice         float64            `gorm:"column:final_price"`
	ExecutorPrice      float64            `gorm:"column:executor_price"`
	ExecutorPriceNote  string             `gorm:"column:executor_price_note"`
	Executor2Price     float64            `gorm:"column:executor2_price"`
	Executor2PriceNote string             `gorm:"column:executor2_price_note"`
	Executor3Price     float64            `gorm:"column:executor3_price"`
	Executor3PriceNote string             `gorm:"column:executor3_price_note"`
	PurchaseDate       *time.Time         `gorm:"column:purchase_date"`
	CompletionDate     *time.Time         `gorm:"column:completion_date"`
	StartImages        string             `gorm:"start_images"`         // 订单开始时上传的图片
	DetailImages       xmysql.StringArray `gorm:"detail_images"`        // 订单结束后上传的图片
	CutRate            float64            `gorm:"column:cut_rate"`      // 抽成比例
	IsEvaluation       bool               `gorm:"column:is_evaluation"` // 是否完成评价
	OutRefundNo        string             `gorm:"column:out_refund_no"`
	Snapshot           string             `gorm:"column:snapshot"` // 订单完成时候的快照
	GrabAt             *time.Time         `gorm:"column:grab_at"`  // 打手抢单时间点
}

// TableName sets the table name for the Order model.
func (*Order) TableName() string {
	return "orders"
}

func (o *Order) FetchGameId() string {
	if o == nil || o.RoleId == "" {
		return ""
	}
	id, err := ExtractID(o.RoleId)
	if err != nil {
		return ""
	}
	return id
}

func (o *Order) FetchRoleId() string {
	id, err := ExtractRole(o.RoleId)
	if err != nil {
		return o.RoleId
	}
	return id
}

var (
	// 定义正则表达式
	gameIdRegex = regexp.MustCompile(`Id:(\d+)`)

	// 定义正则表达式
	roleRegex = regexp.MustCompile(`角色:\s*([\p{Han}\w]+)`)
)

// ExtractID 提取 Id: 后的数字部分, Id:9913193068,角色:六哥会打呀 => 9913193068
func ExtractID(input string) (string, error) {
	// 匹配结果
	match := gameIdRegex.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1], nil // 返回捕获组中的数字
	}
	return "", fmt.Errorf("未找到 ID")
}

// ExtractRole 提取 角色: 后的中文部分, Id:9913193068,角色:六哥会打呀 => 六哥会打呀
func ExtractRole(input string) (string, error) {
	// 匹配结果
	match := roleRegex.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1], nil // 返回捕获组中的中文
	}
	return "", fmt.Errorf("未找到角色名称")
}
