package penalty

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"math"
	"mxclub/pkg/utils"
	"time"
)

func init() {
	OfferPenaltyStrategy(DeductRuleTimeout, new(TimeoutPenalty))
}

// TimeoutPenalty 实现超时罚款规则
type TimeoutPenalty struct{}

// 定义超时与罚款金额的映射
var penaltyMap = utils.NewLinkedHashMapWithPairs[time.Duration, float64]([]*maps.Pair[time.Duration, float64]{
	{time.Duration(0), 0},
	{time.Minute * 20, 10}, // 超时10分钟罚款10
	{time.Duration(math.MaxInt64), 0},
})

func (p *TimeoutPenalty) ApplyPenalty(req *PenaltyReq) (*PenaltyResp, error) {
	if req == nil || req.GrabTime == nil {
		return nil, errors.New("param is bad")
	}
	// 计算超时分钟数
	minutes := time.Since(*req.GrabTime)
	durations := penaltyMap.KeySet()
	// 遍历映射并找到合适的罚款金额
	for index, duration := range durations {
		if minutes >= duration && minutes < durations[index+1] {
			penalty := penaltyMap.MustGet(duration)
			return &PenaltyResp{
				DeductType:    DeductRuleTimeout,
				PenaltyAmount: penalty,
				Reason: fmt.Sprintf("订单:[%v], 订单原价为: %v, 接单后，%v分钟还没组队完成开始订单，罚款：%v元",
					req.OrdersId, req.OrderRawPrice, duration.Minutes(), penalty),
				Message: fmt.Sprintf(
					"尊敬的打手您好，您的订单: %v, "+
						"由于接单后没有及时组队并开始订单，组队时间为：%v，罚款：%v元, "+
						"您可在两个工作日内向客服发起申述，超过两个工作日，系统将进行罚款",
					req.OrdersId, utils.FormatDuration(minutes), penalty),
			}, nil
		}
	}
	if defaultResp.PenaltyAmount <= 0 {
		return nil, errors.New("no penalty")
	}
	return defaultResp, nil
}

func (p *TimeoutPenalty) MustApplyPenalty(req *PenaltyReq) *PenaltyResp {
	if req == nil || req.GrabTime == nil {
		return defaultResp
	}
	// 计算超时分钟数
	minutes := time.Since(*req.GrabTime)
	durations := penaltyMap.KeySet()
	// 遍历映射并找到合适的罚款金额
	for index, duration := range durations {
		if minutes >= duration && minutes < durations[index+1] {
			return &PenaltyResp{PenaltyAmount: penaltyMap.MustGet(duration)}
		}
	}
	return defaultResp
}
