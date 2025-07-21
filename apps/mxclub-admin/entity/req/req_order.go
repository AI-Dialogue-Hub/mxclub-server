package req

import (
	"mxclub/pkg/api"
)

type OrderListReq struct {
	*api.PageParams
	OrderStatus string `json:"status" validate:"required"`
	Ge          string `json:"search_GE_createTime"` // start time
	Le          string `json:"search_LE_createTime"` // end time
	ExecutorId  int    `json:"executor_id"`
	OrderId     string `json:"order_id"`
}

type WitchDrawListReq struct {
	*api.PageParams
	WithdrawalStatus string `json:"withdrawal_status"`
	DasherId         int    `json:"dasher_id"`
}

type WitchDrawUpdateReq struct {
	Id               uint    `json:"id"`
	WithdrawalStatus string  `json:"withdrawal_status"`
	WithdrawalInfo   string  `json:"withdrawal_info"`                       // 同意提现，提现方式 拒绝提现拒绝原因
	DasherId         int     `json:"dasher_id" validate:"required"`         // 打手编号
	WithdrawalAmount float64 `json:"withdrawal_amount" validate:"required"` // 提现金额
}

type WxPayRefundsReq struct {
	OutTradeNo string `json:"out_trade_no,omitempty" validate:"required" reg_err_info:"不能为空"`
	Reason     string `json:"reason"`
}

type OrderTradeExportReq struct {
	StartDate string `json:"start_date,omitempty" validate:"required" reg_err_info:"不能为空"`
	EndDate   string `json:"end_date"`
}

type HistoryWithDrawAmountReq struct {
	UserId uint `json:"user_id" validate:"required"`
}

type DeactivateReq struct {
	*api.PageParams
	DasherId   int    `json:"dasher_id"`
	DasherName string `json:"dasher_name"`
}
