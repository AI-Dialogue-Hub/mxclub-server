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
	FindByOrderList(ctx jet.Ctx, orderList []uint64) (map[uint64][]*po.OrderEvaluation, error)
	FindByOrderId(ctx jet.Ctx, orderId any, dasherId any) (*po.OrderEvaluation, error)
	RemoveByOrderIdAndDasherId(ctx jet.Ctx, orderId any, dasherId any) error
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

	sql := `SELECT COALESCE(SUM(rating), 0) AS total_score, COUNT(rating) AS total_count 
			FROM order_evaluation 
			WHERE executor_id = ? and deleted_at is null`

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

func (repo EvaluationRepo) FindByOrderList(ctx jet.Ctx, orderList []uint64) (map[uint64][]*po.OrderEvaluation, error) {
	queryWrapper := xmysql.NewMysqlQuery()
	queryWrapper.SetFilter("order_id in (?)", orderList)
	evaluationPOList, err := repo.ListNoCountByQuery(queryWrapper)
	if err != nil {
		ctx.Logger().Errorf("cannot find evaluation by orderIds => %v", orderList)
		return nil, err
	}
	orderId2EvaluationPOMap := make(map[uint64][]*po.OrderEvaluation)
	for _, evaluationPO := range evaluationPOList {
		var orderId = evaluationPO.OrderID
		if _, ok := orderId2EvaluationPOMap[orderId]; ok {
			evaluations, _ := orderId2EvaluationPOMap[orderId]
			evaluations = append(evaluationPOList, evaluationPO)
			orderId2EvaluationPOMap[orderId] = evaluations
		} else {
			evaluations := make([]*po.OrderEvaluation, 1)
			evaluations = append(evaluations, evaluationPO)
			orderId2EvaluationPOMap[orderId] = evaluations
		}
	}
	return orderId2EvaluationPOMap, nil
}

func (repo EvaluationRepo) FindByOrderId(ctx jet.Ctx, orderId any, dasherId any) (*po.OrderEvaluation, error) {
	return repo.FindOne("order_id = ? and executor_id = ?", orderId, dasherId)
}
func (repo EvaluationRepo) RemoveByOrderIdAndDasherId(ctx jet.Ctx, orderId any, dasherId any) error {
	err := repo.Remove("order_id = ? and executor_id = ?", orderId, dasherId)
	if err != nil {
		ctx.Logger().Errorf("[EvaluationRepo#RemoveByOrderIdAndDasherId]ERROR,%v", err)
		return err
	}
	ctx.Logger().Infof("[EvaluationRepo#RemoveByOrderIdAndDasherId]remove success, orderId:%v, dasherId:%v", orderId, dasherId)
	return nil
}
