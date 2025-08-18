package req

type WxPayReq struct {
	Amount float64 `json:"amount" validate:"required" reg_err_info:"不能为空"`
}

type WxPayRefundsReq struct {
	OrderId string `json:"order_id,omitempty" validate:"required" reg_err_info:"不能为空"`
	Reason  string `json:"reason"`
}
