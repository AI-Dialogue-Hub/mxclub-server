package req

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/api"
)

type OrderListReq struct {
	api.PageParams
	OrderStatus enum.OrderStatus `json:"status"`
	Ge          string           `json:"search_GE_createTime"` // start time
	Le          string           `json:"search_LE_createTime"` // end time
}

type OrderReq struct {
	ExecutorId      uint   `json:"executor_id" validate:"required"`
	GameRegion      string `json:"game_region" validate:"required"`
	Notes           string `json:"notes" validate:"required"`
	ProductId       uint   `json:"product_id" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	SpecifyExecutor bool   `json:"specify_executor" validate:"required"`
	RoleId          string `json:"role_id" validate:"required"`
	OrderName       string `json:"order_name"`
	OrderIcon       string `json:"order_icon"`
}
