package enum

type MiniConfigEnum string

var (
	Swiper           MiniConfigEnum = "swiper"
	Notifications    MiniConfigEnum = "notifications"         // 通知栏
	MyMessage        MiniConfigEnum = "mymessage"             // 我的消息
	ProductType      MiniConfigEnum = "product_type"          // 我的消息
	CutRate          MiniConfigEnum = "cut_rate"              // 抽成比例
	DelayTime        MiniConfigEnum = "gold_order_delay_time" // 非金牌单打手多久能看到单
	DasherEvaluation MiniConfigEnum = "dasher_evaluation"     // 打手评价
	WarningInfo      MiniConfigEnum = "warningInfo"           // warningInfo
)

var miniConfigEnumMap = map[MiniConfigEnum]string{
	Swiper:           "轮播图",
	Notifications:    "通知栏",
	MyMessage:        "我的消息",
	ProductType:      "商品类型",
	CutRate:          "抽成比例(百分之多少，填整数，例如20)",
	DelayTime:        "金牌打手可以提前多少秒看到单(例如：20 单位s)",
	DasherEvaluation: "提示给用户看到评价信息",
	WarningInfo:      "提示打手接单的信息",
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
