package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/lottery/repo"
)

func init() {
	jet.Provide(NewLotteryService)
}

type LotteryService struct {
	lotteryPrizeRepo    repo.ILotteryPrizeRepo
	lotteryActivityRepo repo.ILotteryActivityRepo
	lotteryRepo         repo.ILotteryRepo
}

func NewLotteryService(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivityRepo repo.ILotteryActivityRepo,
	lotteryRepo repo.ILotteryRepo) *LotteryService {
	return &LotteryService{
		lotteryPrizeRepo:    lotteryPrizeRepo,
		lotteryActivityRepo: lotteryActivityRepo,
		lotteryRepo:         lotteryRepo,
	}
}
