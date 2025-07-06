package dto

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"mxclub/domain/lottery/entity/enum"
)

type LotteryPrizeTypeDTO struct {
	PrizeType maps.IMap[enum.PrizeTypeEnum, string] `json:"prize_type"`
}
