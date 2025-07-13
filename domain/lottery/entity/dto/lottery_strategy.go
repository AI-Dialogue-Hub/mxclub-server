package dto

import (
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
)

type LotteryStrategyBeforeDrawDTO struct {
	UserId                       uint
	ActivityId                   uint
	PrizeId                      uint
	PrizeLevel                   enum.LotteryPrizeLevelEnum
	ActivityPrizeSnapshot        string
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
	UserId                       uint
	ActivityId                   uint
	PrizeId                      uint
	PrizeLevel                   enum.LotteryPrizeLevelEnum
	ActivityPrizeSnapshot        string
	IsIncreaseLotteryProbability bool // 是否提升抽奖概率
}

func (l *LotteryStrategyDrawDTO) WrapBeforeDrawDTO() *LotteryStrategyBeforeDrawDTO {
	return &LotteryStrategyBeforeDrawDTO{
		UserId:                       l.UserId,
		ActivityId:                   l.ActivityId,
		PrizeId:                      l.PrizeId,
		PrizeLevel:                   l.PrizeLevel,
		ActivityPrizeSnapshot:        l.ActivityPrizeSnapshot,
		IsIncreaseLotteryProbability: l.IsIncreaseLotteryProbability,
	}
}

func (l *LotteryStrategyDrawDTO) WrapAfterDrawDTO() *LotteryStrategyAfterDrawDTO {
	return &LotteryStrategyAfterDrawDTO{
		UserId:                l.UserId,
		ActivityId:            l.ActivityId,
		PrizeId:               l.PrizeId,
		PrizeLevel:            l.PrizeLevel,
		ActivityPrizeSnapshot: l.ActivityPrizeSnapshot,
	}
}

type LotteryStrategyDrawResultDTO struct {
	UserId                       uint
	ActivityId                   uint
	PrizeId                      uint
	LotteryPrize                 *po.LotteryPrize // 中奖的奖品
	IsIncreaseLotteryProbability bool             // 是否提升抽奖概率
}
