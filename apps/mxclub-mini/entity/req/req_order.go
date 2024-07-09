package req

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/api"
)

type OrderReq struct {
}

type OrderListReq struct {
	api.PageParams
	OrderStatus enum.OrderStatus `json:"status"`
	Ge          string           `json:"search_GE_createTime"` // start time
	Le          string           `json:"search_LE_createTime"` // end time
}
