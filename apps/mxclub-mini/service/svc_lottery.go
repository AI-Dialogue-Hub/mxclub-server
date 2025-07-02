package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/lottery/repo"
)

func init() {
	jet.Provide(NewLotteryService)
}

type LotteryService struct {
	lotteryPrizeRepo repo.ILotteryPrizeRepo
}

func NewLotteryService(lotteryPrizeRepo repo.ILotteryPrizeRepo) *LotteryService {
	return &LotteryService{lotteryPrizeRepo: lotteryPrizeRepo}
}
