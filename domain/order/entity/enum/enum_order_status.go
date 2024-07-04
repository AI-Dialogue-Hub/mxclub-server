package enum

type OrderStatus int

const (
	PROCESSING     OrderStatus = iota // 配单中
	RUNNING                           // 进行中
	SUCCESS                           // 已完成
	ORDER_ACCEPTED                    // 已接单
	CANCELLED                         // 已取消
	PAUSED                            // 暂停中
)

var OrderStatusMap = map[OrderStatus]string{
	RUNNING:        "进行中",
	SUCCESS:        "已完成",
	ORDER_ACCEPTED: "已接单",
	CANCELLED:      "已取消",
	PROCESSING:     "配单中",
	PAUSED:         "暂停中",
}

func (order OrderStatus) Valid() bool {
	if _, ok := OrderStatusMap[order]; ok {
		return true
	}
	return false
}
