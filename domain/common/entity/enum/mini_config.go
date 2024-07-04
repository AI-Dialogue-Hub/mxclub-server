package enum

type MiniConfigEnum string

var (
	Swiper        MiniConfigEnum = "swiper"
	Notifications MiniConfigEnum = "notifications" // 通知栏
	MyMessage     MiniConfigEnum = "mymessage"     // 我的消息
)

var miniConfigEnumMap = map[MiniConfigEnum]any{
	Swiper:        nil,
	Notifications: nil,
}

func (m MiniConfigEnum) IsValid() bool {
	_, ok := miniConfigEnumMap[m]
	return ok
}

func (m MiniConfigEnum) IsNotValid() bool {
	return !m.IsValid()
}
