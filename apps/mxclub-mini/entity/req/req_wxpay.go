package req

type WxPayReq struct {
	Amount float64 `json:"amount" validate:"required" reg_err_info:"不能为空"`
}

type WxPayRefundsReq struct {
	OutTradeNo string `json:"out_trade_no,omitempty" validate:"required" reg_err_info:"不能为空"`
	Reason     string `json:"reason"`
}
