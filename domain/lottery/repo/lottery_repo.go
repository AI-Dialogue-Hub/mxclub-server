package repo

import "github.com/fengyuan-liang/jet-web-fasthttp/jet"

func init() {
	jet.Provide(NewLotteryRepo)
}

type ILotteryRepo interface {
}

type LotteryRepo struct {
	lotteryActivityRepo ILotteryActivityRepo
	lotteryPrizeRepo    ILotteryPrizeRepo
}

func NewLotteryRepo(lotteryActivityRepo ILotteryActivityRepo, lotteryPrizeRepo ILotteryPrizeRepo) ILotteryRepo {
	return &LotteryRepo{
		lotteryActivityRepo: lotteryActivityRepo,
		lotteryPrizeRepo:    lotteryPrizeRepo,
	}
}
