package service

import (
	"cmp"
	"errors"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	commonEnum "mxclub/domain/common/entity/enum"
	userEnum "mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
	"slices"
	"sync"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
)

// AllDasherHistoryWithDrawAmount 导出所有打手的订单金额记录（优化版：使用批量查询）
func (svc *OrderService) AllDasherHistoryWithDrawAmount(ctx jet.Ctx) ([]*vo.HistoryWithDrawVO, error) {
	defer utils.TraceElapsed(ctx, "AllDasherHistoryWithDrawAmount")()

	// 0. 查询所有打手（限制最大数量）
	const MaxDasherQueryLimit = 10000
	allDasherList, count, err := svc.userRepo.ListAroundCacheByUserTypeAndDasherId(
		ctx, &api.PageParams{Page: 1, PageSize: MaxDasherQueryLimit}, userEnum.RoleAssistant, -1)
	if err != nil {
		ctx.Logger().Errorf("[OrderService#AllDasherHistoryWithDrawAmount] find dasher failed, err:%v", err)
		return nil, errors.New("打手查询失败")
	}
	ctx.Logger().Infof("find dasher count is %v", count)

	if len(allDasherList) == 0 {
		return []*vo.HistoryWithDrawVO{}, nil
	}

	// 1. 准备批量查询所需的ID列表
	dasherIds := make([]int, 0, len(allDasherList))
	userIds := make([]uint, 0, len(allDasherList))
	dasherNumberToUser := make(map[uint]*po.User) // dasherNumber -> User

	for _, dasher := range allDasherList {
		if dasher.MemberNumber < 0 {
			continue
		}
		dasherIds = append(dasherIds, dasher.MemberNumber)
		userIds = append(userIds, dasher.ID)
		dasherNumberToUser[uint(dasher.MemberNumber)] = dasher
	}

	// 2. 并发批量查询所有数据
	var (
		withdrawnAmounts map[int]float64  // 已提现金额（非拒绝）
		approvedAmounts  map[int]float64  // 已成功提现金额
		orderAmounts     map[int]float64  // 订单可提现金额
		rewardAmounts    map[uint]float64 // 打赏金额
		deductAmounts    map[uint]float64 // 罚款金额
		wg               = new(sync.WaitGroup)
		errWithdrawn     error
		errOrder         error
		errReward        error
		errDeduct        error
	)

	wg.Add(4)

	// 批量查询提现金额
	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		withdrawnAmounts, approvedAmounts, errWithdrawn = svc.withdrawRepo.BatchWithdrawAmountByDasherIds(ctx, dasherIds)
	}()

	// 批量查询订单金额
	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		orderAmounts, errOrder = svc.orderRepo.BatchOrderWithdrawAbleAmount(ctx, dasherIds)
	}()

	// 批量查询打赏金额
	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		rewardAmounts, errReward = svc.rewardRecordRepo.BatchRewardAmountByDasherIds(ctx, userIds)
	}()

	// 批量查询罚款金额
	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		deductAmounts, errDeduct = svc.deductionRepo.BatchTotalDeductByUserIds(ctx, userIds)
	}()

	wg.Wait()

	// 检查错误
	if errWithdrawn != nil {
		ctx.Logger().Errorf("[AllDasherHistoryWithDrawAmount] BatchWithdrawAmountByDasherIds error: %v", errWithdrawn)
		return nil, errors.New("批量查询提现金额失败")
	}
	if errOrder != nil {
		ctx.Logger().Errorf("[AllDasherHistoryWithDrawAmount] BatchOrderWithdrawAbleAmount error: %v", errOrder)
		return nil, errors.New("批量查询订单金额失败")
	}
	if errReward != nil {
		ctx.Logger().Errorf("[AllDasherHistoryWithDrawAmount] BatchRewardAmountByDasherIds error: %v", errReward)
		return nil, errors.New("批量查询打赏金额失败")
	}
	if errDeduct != nil {
		ctx.Logger().Errorf("[AllDasherHistoryWithDrawAmount] BatchTotalDeductByUserIds error: %v", errDeduct)
		return nil, errors.New("批量查询罚款金额失败")
	}

	// 3. 组装结果
	withDrawVOList := make([]*vo.HistoryWithDrawVO, 0, len(allDasherList))
	for _, dasher := range allDasherList {
		if dasher.MemberNumber < 0 {
			continue
		}

		dasherId := dasher.MemberNumber
		dasherNumber := uint(dasherId)

		// 获取各项金额
		approveWithdrawnAmount := utils.RoundToTwoDecimalPlaces(approvedAmounts[dasherId])
		orderWithdrawAbleAmount := utils.RoundToTwoDecimalPlaces(orderAmounts[dasherId])
		rewardAmount := utils.RoundToTwoDecimalPlaces(rewardAmounts[dasher.ID])
		withdrawnAmount := utils.RoundToTwoDecimalPlaces(withdrawnAmounts[dasherId])
		totalDeduct := utils.RoundToTwoDecimalPlaces(deductAmounts[dasher.ID])

		// 计算可提现金额
		withdrawAbleAmount := utils.RoundToTwoDecimalPlaces(
			orderWithdrawAbleAmount + rewardAmount - withdrawnAmount - totalDeduct)

		// 获取提现范围
		minRangeNum, maxRangeNum := svc.fetchWithDrawRange(ctx)

		withDrawVOList = append(withDrawVOList, &vo.HistoryWithDrawVO{
			WithDrawVO: &vo.WithDrawVO{
				HistoryWithDrawAmount: approveWithdrawnAmount,
				WithdrawAbleAmount:    withdrawAbleAmount,
				WithdrawRangeMax:      float64(maxRangeNum),
				WithdrawRangeMin:      float64(minRangeNum),
			},
			DasherID:   dasherNumber,
			DasherName: dasher.Name,
		})
	}

	ctx.Logger().Infof("withDrawVOList count is %v", len(withDrawVOList))

	// 4. 按 DasherID 排序
	slices.SortFunc(withDrawVOList, func(a, b *vo.HistoryWithDrawVO) int {
		return cmp.Compare(a.DasherID, b.DasherID)
	})

	return withDrawVOList, nil
}

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
	defer utils.RecoverAndLogError(ctx)

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
