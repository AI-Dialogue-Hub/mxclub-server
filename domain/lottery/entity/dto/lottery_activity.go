package dto

import "mxclub/domain/lottery/po"

type LotteryActivityDTO struct {
	LotteryActivity *po.LotteryActivity `json:"lottery_activity"`
	LotteryPrizes   []*po.LotteryPrize  `json:"lottery_prize"`
}
