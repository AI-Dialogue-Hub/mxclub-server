package req

type TransferReq struct {
	OrderId      uint64 `json:"order_id"` // 这里的是db id
	ExecutorTo   int    `json:"executor_to"`
	ExecutorFrom int    `json:"executor_from"`
	IsValid      bool   `json:"is_valid"` // ExecutorTo 是否合法 排除0的干扰
}
