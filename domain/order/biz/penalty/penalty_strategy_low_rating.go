package penalty

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"mxclub/pkg/utils"
)

func init() {
	OfferPenaltyStrategy(DeductRuleLowRating, new(LowRatingPenalty))
}

// LowRatingPenalty 实现低评星罚款规则
type LowRatingPenalty struct{}

// 定义超时与罚款金额的映射
var penaltyLowRatingMap = utils.NewLinkedHashMapWithPairs[int, float64]([]*maps.Pair[int, float64]{
	{0, 30},
	{1, 30}, // 超时10分钟罚款10
	{2, 20},
	{3, 0},
})

func (p *LowRatingPenalty) ApplyPenalty(req *PenaltyReq) (*PenaltyResp, error) {
	if req == nil || req.Rating <= 0 {
		return nil, errors.New("param is bad")
	}
	lowRatingSet := penaltyLowRatingMap.KeySet()
	for index, rating := range lowRatingSet {
		if index+1 < len(lowRatingSet) && req.Rating >= rating && req.Rating < lowRatingSet[index+1] {
			penalty := penaltyLowRatingMap.MustGet(rating)
			return &PenaltyResp{
				PenaltyAmount: penalty,
				Reason:        fmt.Sprintf("老板评价低星罚款，评星为：%v星，罚款：%v元", rating, penalty),
				Message: fmt.Sprintf(
					`尊敬的打手您好，您的订单: %v, 由于老板评价:%v星，罚款: %v元，
							您可在五个工作日内向客服发起申述，超过五个工作日，系统将进行罚款`,
					req.OrdersId, rating, penalty),
			}, nil
		}
	}
	// 实现低评星罚款的逻辑
	return defaultResp, nil
}

func (p *LowRatingPenalty) MustApplyPenalty(req *PenaltyReq) *PenaltyResp {
	if req == nil || req.Rating <= 0 {
		return defaultResp
	}
	lowRatingSet := penaltyLowRatingMap.KeySet()
	for index, rating := range lowRatingSet {
		if index+1 < len(lowRatingSet) && req.Rating >= rating && req.Rating < lowRatingSet[index+1] {
			return &PenaltyResp{PenaltyAmount: penaltyLowRatingMap.MustGet(rating)}
		}
	}
	// 实现低评星罚款的逻辑
	return defaultResp
}
