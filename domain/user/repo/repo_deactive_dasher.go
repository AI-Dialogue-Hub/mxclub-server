package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/user/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewDeactivateDasherRepo)
}

type IDeactivateDasherRepo interface {
	xmysql.IBaseRepo[po.DeactivateDasher]
}

func NewDeactivateDasherRepo(db *gorm.DB) IDeactivateDasherRepo {
	repo := new(DeactivateDasherRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.DeactivateDasher)
	repo.Ctx = context.Background()
	return repo
}

type DeactivateDasherRepo struct {
	xmysql.BaseRepo[po.DeactivateDasher]
}
