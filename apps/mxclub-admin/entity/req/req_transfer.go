package req

import (
	"mxclub/pkg/api"
)

type TransferListReq struct {
	api.PageParams
	Status int `form:"status"`
}

type TransferReq struct {
	OrderId    uint64 `json:"order_id" comment:"订单Id"`
	ExecutorTo int    `json:"executor_to" comment:"转单to Id，必须在线 && 不在订单中"`
}
