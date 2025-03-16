package req

import "mxclub/pkg/api"

type RewardListReq struct {
	*api.PageParams
	OrderId int `json:"order_id" form:"order_id"` // 打赏金额
}

type RewardPrepayReq struct {
	RewardAmount float64 `json:"reward_amount"` // 打赏金额
	OrderId      string  `json:"order_id"`      // 订单id
	DasherId     int     `json:"dasher_id"`     // 打手Id
	RewardNote   string  `json:"reward_note"`   // 打赏备注
}
