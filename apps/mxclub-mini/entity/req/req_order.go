package req

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/api"
)

type OrderListReq struct {
	api.PageParams
	OrderStatus  enum.OrderStatus `json:"status"`
	Ge           string           `json:"search_GE_createTime"` // start time
	Le           string           `json:"search_LE_createTime"` // end time
	MemberNumber int              `json:"member_number"`
	Role         string           `json:"role"`
}

type OrderReq struct {
	ExecutorId      int    `json:"executor_id"`
	GameRegion      string `json:"game_region"`
	Notes           string `json:"notes"`
	ProductId       uint   `json:"product_id"`
	Phone           string `json:"phone"`
	SpecifyExecutor bool   `json:"specify_executor"`
	RoleId          string `json:"role_id"`
	OrderName       string `json:"order_name"`
	OrderIcon       string `json:"order_icon"`
	OrderTradeNo    string `json:"out_trade_no"`
}

type OrderFinishReq struct {
	OrderId uint     `json:"order_id" validate:"required"` // 订单db主键
	Images  []string `json:"images" validate:"required"`
}

type OrderStartReq struct {
	OrderId     uint   `json:"orderId" validate:"required"` // 订单流水号
	ExecutorId  int    `json:"executor_id" validate:"required"`
	Executor2Id int    `json:"executor_2_id"`
	Executor3Id int    `json:"executor_3_id"`
	RoleId      string `json:"role_id" validate:"required"`
	GameRegion  string `json:"game_region" validate:"required"`
}

type WithDrawReq struct {
	Amount float64 `json:"amount" validate:"required"`
}

type WithDrawListReq struct {
	*api.PageParams
	Ge string `json:"search_GE_createTime"` // start time
	Le string `json:"search_LE_createTime"` // end time
}

type OrderExecutorReq struct {
	ExecutorType uint   `json:"executor_type" validate:"required"`
	ExecutorName string `json:"executor_name"`
	ExecutorId   int    `json:"executor_id"`
	OrderId      uint   `json:"order_id"`
}

type OrderExecutorInviteReq struct {
	ExecutorId int    `json:"executor_id" validate:"gt=-1" reg_error_info:"打手编号必须大于等于0"`
	OrderId    uint   `json:"orderId" validate:"required"`
	RoleId     string `json:"role_id" validate:"required"`
	GameRegion string `json:"game_region" validate:"required"`
}

type OrderGrabReq struct {
	OrderId    uint `json:"orderId" validate:"required"` // 订单流水号
	ExecutorId int  `json:"executor_id" validate:"required"`
}
