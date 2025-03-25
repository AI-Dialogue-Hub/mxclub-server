package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/biz/penalty"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/utils"
)

func (svc OrderService) AddEvaluation(ctx jet.Ctx, evaluationReq *req.EvaluationReq) error {

	// 防抖
	if err := xredis.DebounceForOneDay(
		fmt.Sprintf("orderId_%v_executorId_%v", evaluationReq.OrderID, evaluationReq.ExecutorID),
	); err != nil {
		ctx.Logger().Errorf("duplicated evaluation,evaluationReq => %v", utils.ObjToJsonStr(evaluationReq))
		return errors.New("已评价成功")
	}

	var evaluationList = make([]*po.OrderEvaluation, 0)

	evaluation1 := &po.OrderEvaluation{
		OrdersID:   evaluationReq.OrdersID,
		OrderID:    uint64(evaluationReq.OrderID),
		ExecutorID: evaluationReq.ExecutorID,
		Rating:     evaluationReq.Rating,
		Comments:   evaluationReq.Comments,
	}

	go svc.handleLowRatingDeduction(ctx, evaluation1)

	evaluationList = append(evaluationList, evaluation1)
	if evaluationReq.Executor2ID >= 0 && evaluationReq.Rating2 > 0 {
		evaluation2 := &po.OrderEvaluation{
			OrdersID:   evaluationReq.OrdersID,
			OrderID:    uint64(evaluationReq.OrderID),
			ExecutorID: evaluationReq.Executor2ID,
			Rating:     evaluationReq.Rating2,
			Comments:   evaluationReq.Comments2,
		}
		evaluationList = append(evaluationList, evaluation2)
		go svc.handleLowRatingDeduction(ctx, evaluation2)
	}
	if evaluationReq.Executor3ID >= 0 && evaluationReq.Rating3 > 0 {
		evaluation3 := &po.OrderEvaluation{
			OrdersID:   evaluationReq.OrdersID,
			OrderID:    uint64(evaluationReq.OrderID),
			ExecutorID: evaluationReq.Executor3ID,
			Rating:     evaluationReq.Rating3,
			Comments:   evaluationReq.Comments3,
		}
		evaluationList = append(evaluationList, evaluation3)
		go svc.handleLowRatingDeduction(ctx, evaluation3)
	}

	many, err := svc.evaluationRepo.InsertMany(evaluationList)
	if err != nil {
		ctx.Logger().Errorf("[AddEvaluation]insertNum:%v, ERROR %v", many, err)
		return errors.New("评价失败")
	}

	// 修改订单的评价状态
	_ = svc.orderRepo.DoneEvaluation(evaluationReq.OrdersID)

	// 低评星进行扣款

	return nil
}

func (svc OrderService) handleLowRatingDeduction(ctx jet.Ctx, evaluation *po.OrderEvaluation) {
	defer utils.RecoverWithPrefix(ctx, "handleLowRatingDeduction")

	var (
		rating  = evaluation.Rating
		orderNo = evaluation.OrderID
		logger  = ctx.Logger()
	)

	if rating > 2 {
		return
	}

	dasherPO, err := svc.userService.FindUserByDashId(ctx, evaluation.ExecutorID)

	if err != nil {
		logger.Errorf("FindUserByDashId ERROR: %v", err)
		return
	}

	penaltyStrategy, err := penalty.FetchPenaltyRule(penalty.DeductRuleLowRating)

	if err != nil {
		logger.Errorf("fetch penaltyRule ERROR: %v", err)
		return
	}

	applyPenalty, err := penaltyStrategy.ApplyPenalty(&penalty.PenaltyReq{OrdersId: uint(orderNo), Rating: rating})

	if err != nil || applyPenalty.PenaltyAmount <= 0 {
		logger.Errorf("fetch penaltyRule ERROR: %v, applyPenalty: %v", err, utils.ObjToJsonStr(applyPenalty))
		return
	}

	var applyPenaltyMessage = fmt.Sprintf("%v,评论:%v", applyPenalty.Reason, evaluation.Comments)

	err = svc.deductionRepo.InsertOne(&po.Deduction{
		UserID:          dasherPO.ID,
		DasherId:        dasherPO.MemberNumber,
		ConfirmPersonId: 0,
		Amount:          applyPenalty.PenaltyAmount,
		Reason:          applyPenaltyMessage,
		Status:          enum.Deduct_PENDING,
		OrderNo:         uint(orderNo),
	})

	_ = svc.messageService.PushSystemMessage(ctx, dasherPO.ID, applyPenaltyMessage)

	if err != nil {
		logger.Errorf("deduction insert ERROR: %v", err)
		return
	}
}

func (svc OrderService) RemoveEvaluation(ctx jet.Ctx) error {
	userId := middleware.MustGetUserId(ctx)
	userPO, _ := svc.userService.FindUserById(ctx, userId)
	return svc.evaluationRepo.RemoveEvaluation(ctx, userPO.MemberNumber)
}
