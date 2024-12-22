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
	FindStaring(ctx jet.Ctx, dasherId int) (float64, error)
	RemoveEvaluation(ctx jet.Ctx, dasherId int) error
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

// ====================================================================================================

func (repo EvaluationRepo) FindStaring(ctx jet.Ctx, dasherId int) (float64, error) {
	type EvaluationStarResult struct {
		TotalScore uint
		TotalCount uint
	}

	var staringResult EvaluationStarResult

	sql := "SELECT COALESCE(SUM(rating), 0) AS total_score, COUNT(rating) AS total_count FROM order_evaluation WHERE executor_id = ?"

	err := repo.DB().
		Raw(sql, dasherId).
		Scan(&staringResult).
		Error
	if err != nil {
		ctx.Logger().Errorf("[FindStaring] ERROR: %v", err)
		return 0, err
	}

	if staringResult.TotalCount == 0 {
		return 0, nil
	}

	averageScore := float64(staringResult.TotalScore) / float64(staringResult.TotalCount)
	return averageScore, nil
}

func (repo EvaluationRepo) RemoveEvaluation(ctx jet.Ctx, dasherId int) error {
	return repo.Remove("executor_id = ?", dasherId)
}
