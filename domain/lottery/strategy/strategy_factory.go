package strategy

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/pkg/errors"
)

var lotteryStrategyFactoryMap = make(map[string]ILotteryStrategy)

func FetchLotteryStrategy(strategyName string) (ILotteryStrategy, error) {
	strategy, ok := lotteryStrategyFactoryMap[strategyName]
	if !ok {
		xlog.Errorf("strategy:%v not found", strategyName)
		return nil, errors.New("strategy not found")
	}
	return strategy, nil
}

func OfferLotteryStrategy(strategyName string, strategy ILotteryStrategy) error {
	if ok := lotteryStrategyFactoryMap[strategyName]; ok != nil {
		xlog.Errorf("strategy:%v already exists", strategyName)
		return errors.New("strategy already exists")
	}
	lotteryStrategyFactoryMap[strategyName] = strategy
	return nil
}
