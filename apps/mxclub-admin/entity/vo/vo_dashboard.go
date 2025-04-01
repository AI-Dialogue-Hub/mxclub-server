package vo

type DashBoardVO struct {
	SaleForDayList []*SaleForDay `json:"sale_for_day_list"`
}

// SaleForDay 指定时间段内的销量
type SaleForDay struct {
	Day       string `json:"day"`  // like 2025-03-30
	OrderSale int    `json:"sale"` // 订单销量
}
