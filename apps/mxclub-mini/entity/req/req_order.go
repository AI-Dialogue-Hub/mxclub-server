package req

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/api"
)

type OrderReq struct {
}

type OrderListReq struct {
	api.PageParams
	OrderStatus enum.OrderStatus `json:"order_status"`
}
