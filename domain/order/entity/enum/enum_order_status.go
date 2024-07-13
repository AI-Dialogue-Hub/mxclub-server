package enum

import "github.com/fengyuan-liang/GoKit/collection/maps"

type OrderStatus int

const (
	ALL            OrderStatus = iota // all
	PROCESSING                        // 配单中
	RUNNING                           // 进行中
	SUCCESS                           // 已完成
	ORDER_ACCEPTED                    // 已接单
	CANCELLED                         // 已取消
	PAUSED                            // 暂停中
)

var OrderStatusMap = map[OrderStatus]string{
	ALL:            "全部订单",
	RUNNING:        "进行中",
	SUCCESS:        "已完成",
	ORDER_ACCEPTED: "已接单",
	CANCELLED:      "已取消",
	PROCESSING:     "配单中",
	PAUSED:         "暂停中",
}

var OrderStatusStringMap = map[string]OrderStatus{
	"ALL":            ALL,
	"RUNNING":        RUNNING,
	"SUCCESS":        SUCCESS,
	"ORDER_ACCEPTED": ORDER_ACCEPTED,
	"CANCELLED":      CANCELLED,
	"PROCESSING":     PROCESSING,
	"PAUSED":         PAUSED,
}

var OrderStatusParseMap = func() maps.IMap[OrderStatus, string] {
	linkedHashMap := maps.NewLinkedHashMap[OrderStatus, string]()
	for k, v := range OrderStatusStringMap {
		linkedHashMap.Put(v, k)
	}
	return linkedHashMap
}()

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
	return OrderStatusParseMap.MustGet(order)
}