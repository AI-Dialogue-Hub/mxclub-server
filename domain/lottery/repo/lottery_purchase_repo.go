package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryPurchaseRecordsRepo)
}

type ILotteryPurchaseRecordsRepo interface {
	xmysql.IBaseRepo[po.LotteryPurchaseRecord]
}

func NewLotteryPurchaseRecordsRepo(db *gorm.DB) ILotteryPurchaseRecordsRepo {
	repo := new(LotteryPurchaseRecordsRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.LotteryPurchaseRecord)
	repo.Ctx = context.Background()
	return repo
}

type LotteryPurchaseRecordsRepo struct {
	xmysql.BaseRepo[po.LotteryPurchaseRecord]
}
