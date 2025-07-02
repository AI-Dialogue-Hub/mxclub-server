package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryPrizeRepo)
}

type ILotteryPrizeRepo interface {
	xmysql.IBaseRepo[po.LotteryPrize]
}

func NewLotteryPrizeRepo(db *gorm.DB) ILotteryPrizeRepo {
	repo := new(LotteryPrizeRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.LotteryPrize)
	repo.Ctx = context.Background()
	return repo
}

type LotteryPrizeRepo struct {
	xmysql.BaseRepo[po.LotteryPrize]
}
