package penalty

import "time"

// =========== 罚款规则和金额 ========================

type (
	DeductRule int
)

const (
	DeductRuleTimeout   DeductRule = iota // 超时
	DeductRuleLowRating                   // 低评星
)

var (
	defaultResp = new(PenaltyResp)
)

// PenaltyReq 定义PenaltyReq和PenaltyResp类型
type PenaltyReq struct {
	OrdersId uint       // 订单Id
	Rating   int        // 评星
	GrabTime *time.Time // 抢单时间
}

type PenaltyResp struct {
	PenaltyAmount float64 // 罚款的金额
	Reason        string  // 罚款原因
	Message       string  // 给打手发的消息
}
