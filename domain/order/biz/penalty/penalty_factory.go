package penalty

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/GoKit/collection/maps"
)

// penaltyStrategyFactory 策略工厂
var penaltyStrategyFactory = maps.NewLinkedHashMap[DeductRule, IPenaltyRule]()

// OfferPenaltyStrategy 添加策略
func OfferPenaltyStrategy(deductRule DeductRule, strategy IPenaltyRule) {
	penaltyStrategyFactory.Put(deductRule, strategy)
}

// FetchPenaltyRule 获取具体处罚方法
func FetchPenaltyRule(rule DeductRule) (IPenaltyRule, error) {
	if penaltyStrategy, ok := penaltyStrategyFactory.Get(rule); ok {
		return penaltyStrategy, nil
	}
	return nil, errors.New(fmt.Sprintf("unknown deduction rule: %v", rule))
}
