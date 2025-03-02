package repo

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewProductSalesRepo)
}

type IRewardRecordRepo interface {
	xmysql.IBaseRepo[po.RewardRecord]
}

func NewProductSalesRepo(db *gorm.DB) IRewardRecordRepo {
	repo := new(IRewardRecordRepoImpl)
	repo.SetDB(db)
	repo.ModelPO = new(po.RewardRecord)
	return repo
}

type IRewardRecordRepoImpl struct {
	xmysql.BaseRepo[po.RewardRecord]
}
