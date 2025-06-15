package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	commonEnum "mxclub/domain/common/entity/enum"
	userEnum "mxclub/domain/user/entity/enum"
	"mxclub/pkg/utils"
	"sync"
)

func (svc *OrderService) HistoryWithDrawAmount(ctx jet.Ctx, req *req.HistoryWithDrawAmountReq) (*vo.WithDrawVO, error) {
	userId := req.UserId
	userPO, err := svc.userRepo.FindByIdAroundCache(ctx, userId)
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, cannot find user:%v", userId)
		return nil, errors.New("cannot find user info")
	}
	if userPO.Role != userEnum.RoleAssistant || userPO.MemberNumber < 0 {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, user not dasher:%v", userId)
		return nil, errors.New("您还不是打手")
	}
	var (
		approveWithdrawnAmount  float64
		withdrawnAmount         float64
		orderWithdrawAbleAmount float64
		totalDeduct             float64
		rewardAmount            float64 // 打赏的钱
		wg                      = new(sync.WaitGroup)
	)

	wg.Add(5)

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		// 提现成功的钱
		approveWithdrawnAmount, _ = svc.withdrawRepo.ApproveWithdrawnAmount(ctx, userPO.MemberNumber)
		// 四舍五入
		approveWithdrawnAmount = utils.RoundToTwoDecimalPlaces(approveWithdrawnAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		// 用户发起提现的钱，包括未提现和提现成功的
		withdrawnAmount, _ = svc.withdrawRepo.WithdrawnAmountNotReject(ctx, userPO.MemberNumber)
		withdrawnAmount = utils.RoundToTwoDecimalPlaces(withdrawnAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		// 订单中能提现的钱
		orderWithdrawAbleAmount, _ = svc.orderRepo.OrderWithdrawAbleAmount(ctx, userPO.MemberNumber)
		orderWithdrawAbleAmount = utils.RoundToTwoDecimalPlaces(orderWithdrawAbleAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		rewardAmount, _ = svc.rewardRecordRepo.AllRewardAmountByDasherId(ctx, userPO.ID)
		rewardAmount = utils.RoundToTwoDecimalPlaces(rewardAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		// 罚款的钱
		totalDeduct, _ = svc.deductionRepo.TotalDeduct(ctx, userId)
		totalDeduct = utils.RoundToTwoDecimalPlaces(totalDeduct)
	}()

	wg.Wait()

	ctx.Logger().Infof(
		"dashId:%v, approveWithdrawnAmount:%v, "+
			"withdrawnAmount:%v, orderWithdrawAbleAmount:%v,totalDeduct:%v, rewardAmount:%v",
		userPO.MemberNumber, approveWithdrawnAmount,
		withdrawnAmount, orderWithdrawAbleAmount, totalDeduct, rewardAmount,
	)

	if approveWithdrawnAmount > orderWithdrawAbleAmount+rewardAmount {
		ctx.Logger().Errorf(
			"[HistoryWithDrawAmount]ERROR, approveWithdrawnAmount: %v gt orderWithdrawAbleAmount:%v",
			approveWithdrawnAmount, orderWithdrawAbleAmount,
		)
		return nil, errors.New("系统查询错误，请联系管理员")
	}
	minRangeNum, maxRangeNum := svc.fetchWithDrawRange(ctx)

	// 能提现的钱
	withdrawAbleAmount := utils.RoundToTwoDecimalPlaces(
		orderWithdrawAbleAmount + rewardAmount - withdrawnAmount - totalDeduct)

	return &vo.WithDrawVO{
		HistoryWithDrawAmount: utils.RoundToTwoDecimalPlaces(approveWithdrawnAmount),
		WithdrawAbleAmount:    withdrawAbleAmount,
		WithdrawRangeMax:      float64(maxRangeNum),
		WithdrawRangeMin:      float64(minRangeNum),
	}, nil
}

func (svc *OrderService) fetchWithDrawRange(ctx jet.Ctx) (int, int) {
	utils.RecoverAndLogError(ctx)

	// 获取抽成比例
	minRange := svc.commonRepo.FindConfigByNameOrDefault(
		ctx,
		commonEnum.WithdrawRangeMin.String(),
		nil,
	)

	maxRange := svc.commonRepo.FindConfigByNameOrDefault(
		ctx,
		commonEnum.WithdrawRangeMax.String(),
		nil,
	)

	var (
		minRangeNum = 200
		maxRangeNum = 2000
	)

	if minRange != nil && len(minRange.Content) > 0 && minRange.Content[0] != nil && minRange.Content[0]["desc"] != nil {
		minRangeNum = utils.SafeParseNumber[int](minRange.Content[0]["desc"])
	}

	if maxRange != nil && len(maxRange.Content) > 0 && maxRange.Content[0] != nil && maxRange.Content[0]["desc"] != nil {
		maxRangeNum = utils.SafeParseNumber[int](maxRange.Content[0]["desc"])
	}
	return minRangeNum, maxRangeNum
}
