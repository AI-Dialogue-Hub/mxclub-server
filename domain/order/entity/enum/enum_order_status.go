package enum

type OrderStatus int

const (
	ALL            OrderStatus = iota // all
	PROCESSING                        // 配单中 一个订单提交后就变成配单状态，如果没有指定打手会投到订单池中
	RUNNING                           // 进行中
	SUCCESS                           // 已完成
	PAUSED                            // 暂停中
	ORDER_ACCEPTED                    // 已接单
	Refunds                           // 已退单
	PrePay                            // 待支付
)

var OrderStatusMap = map[OrderStatus]string{
	ALL:            "全部订单",
	RUNNING:        "进行中",
	SUCCESS:        "已完成",
	ORDER_ACCEPTED: "已接单",
	PROCESSING:     "配单中",
	PAUSED:         "暂停中",
	Refunds:        "已退单",
	PrePay:         "待支付",
}

var OrderStatusStringMap = map[string]OrderStatus{
	"ALL":            ALL,
	"RUNNING":        RUNNING,
	"SUCCESS":        SUCCESS,
	"ORDER_ACCEPTED": ORDER_ACCEPTED,
	"PROCESSING":     PROCESSING,
	"PAUSED":         PAUSED,
	"Refunds":        Refunds,
	"PrePay":         PrePay,
}

func (order OrderStatus) Valid() bool {
	if _, ok := OrderStatusMap[order]; ok {
		return true
	}
	return false
}

func ParseOrderStatusByString(orderStr string) OrderStatus {
	return OrderStatusStringMap[orderStr]
}

func (order OrderStatus) String() string {
	return OrderStatusMap[order]
}
