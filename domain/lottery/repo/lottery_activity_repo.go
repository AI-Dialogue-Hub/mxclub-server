package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryActivityRepo)
}

type ILotteryActivityRepo interface {
	xmysql.IBaseRepo[po.LotteryActivity]
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
