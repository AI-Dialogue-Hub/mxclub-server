package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/domain/order/po"
)

func (svc OrderService) AddEvaluation(ctx jet.Ctx, evaluationReq *req.EvaluationReq) error {
	var evaluationList = make([]*po.OrderEvaluation, 0)

	evaluation1 := &po.OrderEvaluation{
		OrdersID:   evaluationReq.OrdersID,
		OrderID:    evaluationReq.OrderID,
		ExecutorID: evaluationReq.ExecutorID,
		Rating:     evaluationReq.Rating,
		Comments:   evaluationReq.Comments,
	}
	evaluationList = append(evaluationList, evaluation1)
	if evaluationReq.Executor2ID > 0 {
		evaluation2 := &po.OrderEvaluation{
			OrdersID:   evaluationReq.OrdersID,
			OrderID:    evaluationReq.OrderID,
			ExecutorID: evaluationReq.Executor2ID,
			Rating:     evaluationReq.Rating2,
			Comments:   evaluationReq.Comments2,
		}
		evaluationList = append(evaluationList, evaluation2)
	}
	if evaluationReq.Executor3ID > 0 {
		evaluation3 := &po.OrderEvaluation{
			OrdersID:   evaluationReq.OrdersID,
			OrderID:    evaluationReq.OrderID,
			ExecutorID: evaluationReq.Executor3ID,
			Rating:     evaluationReq.Rating3,
			Comments:   evaluationReq.Comments3,
		}
		evaluationList = append(evaluationList, evaluation3)
	}

	many, err := svc.evaluationRepo.InsertMany(evaluationList)
	if err != nil {
		ctx.Logger().Errorf("[AddEvaluation]insrtNum:%v, ERROR %v", many, err)
		return errors.New("评价失败")
	}

	// 修改订单的评价状态
	_ = svc.orderRepo.DoneEvaluation(evaluationReq.OrdersID)

	return nil
}
