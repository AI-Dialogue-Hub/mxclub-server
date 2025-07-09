package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryActivityRepo)
}

type ILotteryActivityRepo interface {
	xmysql.IBaseRepo[po.LotteryActivity]
	UpdateStatus(ctx jet.Ctx, id uint, status enum.ActivityStatusEnum) error
}

func NewLotteryActivityRepo(db *gorm.DB) ILotteryActivityRepo {
	repo := new(LotteryActivityRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.LotteryActivity)
	repo.Ctx = context.Background()
	return repo
}

type LotteryActivityRepo struct {
	xmysql.BaseRepo[po.LotteryActivity]
}

func (repo *LotteryActivityRepo) UpdateStatus(ctx jet.Ctx, id uint, status enum.ActivityStatusEnum) error {
	if err := repo.UpdateById(map[string]interface{}{"activity_status": status}, id); err != nil {
		ctx.Logger().Errorf("update ability status error: %v", err)
		return errors.Wrap(err, "update ability status error")
	}
	return nil
}
