package req

import (
	"mxclub/pkg/api"
)

type OrderListReq struct {
	*api.PageParams
	OrderStatus string `json:"status" validate:"required"`
	Ge          string `json:"search_GE_createTime"` // start time
	Le          string `json:"search_LE_createTime"` // end time
	ExecutorId  uint   `json:"executor_id"`
}
