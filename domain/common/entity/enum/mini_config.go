package enum

type MiniConfigEnum string

var (
	Swiper        MiniConfigEnum = "swiper"
	Notifications MiniConfigEnum = "notifications" // 通知栏
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
