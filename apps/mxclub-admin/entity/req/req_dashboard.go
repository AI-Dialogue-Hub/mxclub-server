package req

type SaleForDayReq struct {
	StartDay string `json:"start_day"`
	EndDay   string `json:"end_day"`
}
