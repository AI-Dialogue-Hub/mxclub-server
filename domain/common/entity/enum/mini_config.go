package enum

type MiniConfigEnum string

var (
	Swiper        MiniConfigEnum = "swiper"
	Notifications MiniConfigEnum = "notifications" // 通知栏
	MyMessage     MiniConfigEnum = "mymessage"     // 我的消息
	ProductType   MiniConfigEnum = "product_type"  // 我的消息
)

var miniConfigEnumMap = map[MiniConfigEnum]string{
	Swiper:        "轮播图",
	Notifications: "通知栏",
	MyMessage:     "我的消息",
	ProductType:   "商品类型",
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
