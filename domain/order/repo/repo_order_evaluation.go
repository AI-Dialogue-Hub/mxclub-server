package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewEvaluationRepo)
}

type IEvaluationRepo interface {
	xmysql.IBaseRepo[po.OrderEvaluation]
}

func NewEvaluationRepo(db *gorm.DB) IEvaluationRepo {
	repo := new(EvaluationRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.OrderEvaluation)
	repo.Ctx = context.Background()
	return repo
}

type EvaluationRepo struct {
	xmysql.BaseRepo[po.OrderEvaluation]
}
