package enum

type MiniConfigEnum string

var (
	Swiper              MiniConfigEnum = "swiper"
	Notifications       MiniConfigEnum = "notifications"          // 通知栏
	MyMessage           MiniConfigEnum = "mymessage"              // 我的消息
	ProductType         MiniConfigEnum = "product_type"           // 我的消息
	CutRate             MiniConfigEnum = "cut_rate"               // 抽成比例
	DelayTime           MiniConfigEnum = "gold_order_delay_time"  // 非金牌单打手多久能看到单
	DasherEvaluation    MiniConfigEnum = "dasher_evaluation"      // 打手评价
	WarningInfo         MiniConfigEnum = "warningInfo"            // warningInfo
	SellingPoint        MiniConfigEnum = "SellingPoint"           // 滚动卖点
	WithdrawRangeMin    MiniConfigEnum = "WithdrawRangeMin"       // 滚动卖点
	WithdrawRangeMax    MiniConfigEnum = "WithdrawRangeMax"       // 滚动卖点
	SalesThreshold      MiniConfigEnum = "sales_threshold"        // 出销量热门标阈值
	PayOrderWarningInfo MiniConfigEnum = "pay_order_warning_info" // 订单页面warning信息
	OrderReward         MiniConfigEnum = "order_reward"           // 打赏金额
	payOrderFaker       MiniConfigEnum = "pay_order_faker"        // 打赏金额
)

var miniConfigEnumMap = map[MiniConfigEnum]string{
	Swiper:              "轮播图",
	Notifications:       "通知栏",
	MyMessage:           "我的消息",
	ProductType:         "商品类型",
	CutRate:             "抽成比例(百分之多少，填整数，例如20)",
	DelayTime:           "金牌打手可以提前多少秒看到单(例如：20 单位s)",
	DasherEvaluation:    "提示给用户看到评价信息",
	WarningInfo:         "提示打手接单的信息",
	SellingPoint:        "滚动卖点",
	WithdrawRangeMin:    "最小提现阈值",
	WithdrawRangeMax:    "最大提现阈值",
	SalesThreshold:      "出销量热门标阈值",
	PayOrderWarningInfo: "订单页面warning信息",
	OrderReward:         "打赏金额(直接输入数字不要带单位)",
	payOrderFaker:       "假下单页面(1表示开启，其他为关闭)",
}

func (m MiniConfigEnum) IsValid() bool {
	_, ok := miniConfigEnumMap[m]
	return ok
}

func (m MiniConfigEnum) IsNotValid() bool {
	return !m.IsValid()
}

func (m MiniConfigEnum) DisPlayName() string {
	displayName, ok := miniConfigEnumMap[m]
	if !ok {
		return ""
	}
	return displayName
}

func (m MiniConfigEnum) String() string {
	return string(m)
}
