package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryRecordsRepo)
}

type ILotteryRecordsRepo interface {
	xmysql.IBaseRepo[po.LotteryRecords]
}

func NewLotteryRecordsRepo(db *gorm.DB) ILotteryRecordsRepo {
	repo := new(LotteryRecordsRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.LotteryRecords)
	repo.Ctx = context.Background()
	return repo
}

type LotteryRecordsRepo struct {
	xmysql.BaseRepo[po.LotteryRecords]
}
