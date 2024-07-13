package req

type WxPayReq struct {
	Amount float64 `json:"amount" validate:"required" reg_err_info:"不能为空"`
}
