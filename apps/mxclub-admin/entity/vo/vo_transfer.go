package vo

import (
	"mxclub/domain/order/entity/enum"
	"time"
)

type TransferVO struct {
	ID               uint              `json:"id"`
	OrderId          uint64            `json:"order_id" comment:"订单Id"`
	ExecutorFrom     int               `json:"executor_from" comment:"转单人Id，一般是车头"`
	ExecutorFromName string            `json:"executor_from_name"`
	ExecutorTo       int               `json:"executor_to" comment:"转单to Id，必须在线 && 不在订单中"`
	ExecutorToName   string            `json:"executor_to_name"`
	Status           enum.TransferEnum `json:"status" comment:"转单状态"`
	CreatedAt        time.Time         `json:"created_at"`
}
