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

type WitchDrawListReq struct {
	*api.PageParams
	WithdrawalStatus string `json:"withdrawal_status"`
}

type WitchDrawUpdateReq struct {
	Id                   uint   `json:"id"`
	WithdrawalStatus     string `json:"withdrawal_status"`
	WithdrawalMethod     string `json:"withdrawal_method"`                     // 同意提现，提现方式
	WithdrawalRejectInfo string `json:"withdrawal_reject_info"`                // 拒绝原因
	DasherId             int    `json:"dasher_id" validate:"required"`         // 打手编号
	WithdrawalAmount     string `json:"withdrawal_amount" validate:"required"` // 提现金额
}

type WxPayRefundsReq struct {
	OutTradeNo string `json:"out_trade_no,omitempty" validate:"required" reg_err_info:"不能为空"`
	Reason     string `json:"reason"`
}
