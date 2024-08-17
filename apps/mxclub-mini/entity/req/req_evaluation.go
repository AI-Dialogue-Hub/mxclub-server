package req

type EvaluationReq struct {
	OrdersID    uint   `json:"orders_id"`
	OrderID     uint   `json:"order_id"`
	ExecutorID  int    `json:"executor_id"`
	Executor2ID int    `json:"executor_2_id"`
	Executor3ID int    `json:"executor_3_id"`
	Rating      int    `json:"rating"`
	Comments    string `json:"comments"`
	Rating2     int    `json:"rating2"`
	Comments2   string `json:"comments2"`
	Rating3     int    `json:"rating3"`
	Comments3   string `json:"comments3"`
}
