package dto

import (
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
)

type LotteryStrategyBeforeDrawDTO struct {
	UserId                       uint
	ActivityId                   uint
	IsIncreaseLotteryProbability bool // 是否提升抽奖概率
}

type LotteryStrategyAfterDrawDTO struct {
	UserId                uint
	ActivityId            uint
	PrizeId               uint
	PrizeLevel            enum.LotteryPrizeLevelEnum
	ActivityPrizeSnapshot string
}

type LotteryStrategyDrawDTO struct {
	UserId     uint
	ActivityId uint
}

func (l *LotteryStrategyDrawDTO) WrapBeforeDrawDTO() *LotteryStrategyBeforeDrawDTO {
	return &LotteryStrategyBeforeDrawDTO{
		UserId:     l.UserId,
		ActivityId: l.ActivityId,
	}
}

func (l *LotteryStrategyDrawDTO) WrapAfterDrawDTO() *LotteryStrategyAfterDrawDTO {
	return &LotteryStrategyAfterDrawDTO{
		UserId:     l.UserId,
		ActivityId: l.ActivityId,
	}
}

type LotteryStrategyDrawResultDTO struct {
	PrizeId                      uint
	PrizeIndex                   int              // 奖品在奖池中的位置，默认是优先级
	LotteryPrize                 *po.LotteryPrize // 中奖的奖品
	ActivityPrizeSnapshot        string
	IsIncreaseLotteryProbability bool // 是否提升抽奖概率
}
