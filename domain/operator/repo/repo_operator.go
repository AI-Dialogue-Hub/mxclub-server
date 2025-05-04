package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/operator/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewOperatorLogRepo)
}

type IOperatorRepo interface {
	xmysql.IBaseRepo[po.OperatorLogPO]
}

func NewOperatorLogRepo(db *gorm.DB) IOperatorRepo {
	repo := new(OperatorLogRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.OperatorLogPO)
	repo.Ctx = context.Background()
	return repo
}

type OperatorLogRepo struct {
	xmysql.BaseRepo[po.OperatorLogPO]
}
