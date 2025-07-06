package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryActivityPrizeRelationRepo)
}

type ILotteryActivityPrizeRelationRepo interface {
	xmysql.IBaseRepo[po.ActivityPrizeRelation]
	DelByPrizeId(ctx jet.Ctx, prizeId uint) error
	DelByActivityId(ctx jet.Ctx, activityId uint) error
}

func NewLotteryActivityPrizeRelationRepo(db *gorm.DB) ILotteryActivityPrizeRelationRepo {
	repo := new(LotteryActivityPrizeRelationRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.ActivityPrizeRelation)
	repo.Ctx = context.Background()
	return repo
}

type LotteryActivityPrizeRelationRepo struct {
	xmysql.BaseRepo[po.ActivityPrizeRelation]
}

func (repo *LotteryActivityPrizeRelationRepo) DelByPrizeId(ctx jet.Ctx, prizeId uint) error {
	return repo.Remove("prize_id = ?", prizeId)
}

func (repo *LotteryActivityPrizeRelationRepo) DelByActivityId(ctx jet.Ctx, activityId uint) error {
	return repo.Remove("activity_id = ?", activityId)
}
