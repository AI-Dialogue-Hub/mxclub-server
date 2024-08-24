package req

type TransferReq struct {
	OrderId    uint64 `json:"order_id"`
	ExecutorTo int    `json:"executor_to"`
	IsValid    bool   `json:"is_valid"` // ExecutorTo 是否合法 排除0的干扰
}
