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
	RemoveByPrizeIds(ctx jet.Ctx, ids []uint) (delCount int64, err error)
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

func (repo *LotteryPrizeRepo) RemoveByPrizeIds(ctx jet.Ctx, ids []uint) (delCount int64, err error) {
	tx := repo.DB().Where("id in (?)", ids).Delete(repo.ModelPO)
	tx.Scan(&delCount)
	return delCount, tx.Error
}
