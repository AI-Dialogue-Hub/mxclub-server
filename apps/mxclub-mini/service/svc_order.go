package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	traceUtil "github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"math"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/biz/penalty"
	orderDTO "mxclub/domain/order/entity/dto"
	userPOInfo "mxclub/domain/user/po"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/constant"
	"sync"

	commonEnum "mxclub/domain/common/entity/enum"
	commonRepo "mxclub/domain/common/repo"
	orderRepoDTO "mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	userEnum "mxclub/domain/user/entity/enum"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
	"time"
)

func init() {
	jet.Provide(NewOrderService)
}

type OrderService struct {
	orderRepo      repo.IOrderRepo
	withdrawalRepo repo.IWithdrawalRepo
	userService    *UserService
	productService *ProductService
	messageService *MessageService
	commonRepo     commonRepo.IMiniConfigRepo
	deductionRepo  repo.IDeductionRepo
	evaluationRepo repo.IEvaluationRepo
	transferRepo   repo.ITransferRepo
}

var orderService *OrderService

func NewOrderService(
	repo repo.IOrderRepo,
	withdrawalRepo repo.IWithdrawalRepo,
	userService *UserService,
	productService *ProductService,
	messageService *MessageService,
	commonRepo commonRepo.IMiniConfigRepo,
	deductionRepo repo.IDeductionRepo,
	evaluationRepo repo.IEvaluationRepo,
	transferRepo repo.ITransferRepo) *OrderService {
	orderService = &OrderService{
		orderRepo:      repo,
		withdrawalRepo: withdrawalRepo,
		userService:    userService,
		productService: productService,
		messageService: messageService,
		commonRepo:     commonRepo,
		deductionRepo:  deductionRepo,
		evaluationRepo: evaluationRepo,
		transferRepo:   transferRepo,
	}
	return orderService
}

// ===============================================================

func (svc OrderService) Add(ctx jet.Ctx, req *req.OrderReq) error {
	if _, err := svc.AddByOrderStatus(ctx, req, enum.PrePay); err != nil {
		return errors.New("订单已经添加成功")
	}
	return nil
}

func (svc OrderService) PaySuccessOrder(ctx jet.Ctx, orderNo uint64) error {
	defer utils.RecoverAndLogError(ctx)
	if orderNo <= 0 {
		ctx.Logger().Errorf("invalid orderNo: %v", orderNo)
		return errors.New("invalid orderNo")
	}
	err := svc.orderRepo.UpdateOrderStatus(ctx, orderNo, enum.PROCESSING)
	if err != nil {
		return errors.New("更新失败")
	}
	orderPO, _ := svc.orderRepo.FindByOrderOrOrdersId(ctx, uint(orderNo))
	var (
		executorId   int
		dasherPO     *userPOInfo.User
		orderTradeNo = utils.SafeParseUint64(orderPO.OrderId)
	)
	// 4. 如果指定订单，给打手发送接单消息
	if orderPO.SpecifyExecutor {
		// 指定打手需要该打手同意
		// 特殊编号打手
		executorId = utils.ParseInt(orderPO.OutRefundNo)
		if executorId == -1 {
			executorId = 0
		}
		dasherPO, err = svc.userService.FindUserByDashId(ctx, executorId)
		if err == nil && dasherPO != nil && dasherPO.ID > 0 {
			go func() {
				defer utils.RecoverAndLogError(ctx)
				ctx.Logger().Infof("指定订单:  order  = %v", orderPO)
				// 发送派单信息
				svc.messageService.PushMessage(
					ctx,
					dto.NewDispatchMessageWithFinalPrice(
						dasherPO.ID, uint(orderTradeNo), orderPO.GameRegion, orderPO.RoleId, "",
						utils.RoundToTwoDecimalPlaces(orderPO.FinalPrice)),
				)
			}()
		}
	} else {
		executorId = -1
	}
	ctx.Logger().Infof("pay success, order: %v", utils.ObjToJsonStr(orderPO))
	return nil
}

func (svc OrderService) AddByOrderStatus(ctx jet.Ctx, req *req.OrderReq, status enum.OrderStatus) (*po.Order, error) {
	var (
		userId = middleware.MustGetUserId(ctx)
		logger = ctx.Logger()
	)
	if req.OrderTradeNo == "" {
		req.OrderTradeNo = wxpay.GenerateUniqueOrderNumber()
	}
	// 检查订单是否已经创建
	order, err := svc.orderRepo.FindByOrderId(ctx, utils.ParseUint(req.OrderTradeNo))

	if order != nil && order.ID > 0 {
		logger.Errorf("has duplicates order, %+v", order)
		return nil, err
	}

	if req.Phone != "" {
		go func() { _ = svc.userService.userRepo.UpdateUserPhone(ctx, userId, req.Phone) }()
	}
	var (
		orderTradeNo = utils.SafeParseUint64(req.OrderTradeNo)
		executorId   int
	)
	if req.SpecifyExecutor {
		executorId = req.ExecutorId
		if executorId == -1 {
			executorId = 0
		}
	}
	// 1.2 折扣信息
	preferentialVO, err := svc.Preferential(ctx, req.ProductId)
	if err != nil {
		logger.Errorf("Preferential ERROR:%v", err)
	}
	// 2. 创建订单
	order = &po.Order{
		OrderId:         orderTradeNo,
		PurchaseId:      userId,
		OrderName:       req.OrderName,
		OrderIcon:       req.OrderIcon,
		OrderStatus:     status,
		OriginalPrice:   preferentialVO.OriginalPrice,
		ProductID:       req.ProductId,
		Phone:           req.Phone,
		GameRegion:      req.GameRegion,
		RoleId:          req.RoleId,
		SpecifyExecutor: req.SpecifyExecutor,
		ExecutorID:      -1,
		Executor2Id:     -1,
		Executor3Id:     -1,
		ExecutorName:    "",
		Notes:           req.Notes,
		DiscountPrice:   preferentialVO.OriginalPrice - preferentialVO.DiscountedPrice,
		FinalPrice:      preferentialVO.DiscountedPrice,
		ExecutorPrice:   0,
		PurchaseDate:    utils.Ptr(time.Now()),
		GrabAt:          nil,
		OutRefundNo:     utils.ParseString(executorId), // 保存下用户选择的打手
	}
	// 3. 保存订单
	err = svc.orderRepo.InsertOne(order)
	if err != nil {
		ctx.Logger().Errorf("[orderService]AddDeduction ERROR, %v", err.Error())
		ctx.Logger().Errorf("order:%v", utils.ObjToJsonStr(order))
		return nil, errors.New("订单保存保存失败，请联系客服")
	}

	return order, nil
}

func (svc OrderService) List(ctx jet.Ctx, req *req.OrderListReq) (*api.PageResult, error) {
	userId := middleware.MustGetUserId(ctx)
	// 特殊编号用户
	var executorName string
	var isDasher = req.Role == string(userEnum.RoleAssistant)
	if isDasher {
		userPO, _ := svc.userService.FindUserByDashId(ctx, req.MemberNumber)
		executorName = userPO.Name
	}
	list, err := svc.orderRepo.ListByOrderStatus(ctx, &orderRepoDTO.ListByOrderStatusDTO{
		Status:       req.OrderStatus,
		PageParams:   utils.Ptr(req.PageParams),
		Ge:           req.Ge,
		Le:           req.Le,
		MemberNumber: req.MemberNumber,
		UserId:       userId,
		IsDasher:     isDasher,
		ExecutorName: executorName,
	})
	if err != nil {
		ctx.Logger().Errorf("[orderService]List ERROR, %v", err.Error())
		return nil, errors.New("查询不到数据")
	}
	orderVOS := utils.CopySlice[*po.Order, *vo.OrderVO](list)
	// 打手查看老板等级
	if req.MemberNumber >= 0 {
		// 获取老板等级
		svc.doBuildUserGrade(ctx, orderVOS)
	}
	return api.WrapPageResult(&req.PageParams, orderVOS, 0), err
}

func (svc OrderService) doBuildUserGrade(ctx jet.Ctx, vos []*vo.OrderVO) {
	// 1. 获取所有用户id
	userIdList := utils.Map[*vo.OrderVO, uint](vos, func(in *vo.OrderVO) uint { return in.PurchaseId })
	if len(userIdList) == 0 {
		return
	}
	// 2. 查询userId对应老板的等级
	userId2GradeMap, err := svc.userService.userRepo.FindGradeByUserIdList(userIdList)
	if err != nil {
		ctx.Logger().Errorf("[doBuildUserGrade]ERROR, %v", err)
		return
	}
	utils.ForEach(vos, func(ele *vo.OrderVO) {
		ele.UserGrade = userId2GradeMap.MustGet(ele.PurchaseId)
	})
}

func (svc OrderService) Preferential(ctx jet.Ctx, productId uint) (result *vo.PreferentialVO, err error) {
	userId := middleware.MustGetUserId(ctx)
	userById, err := svc.userService.FindUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	productVO, err := svc.productService.FindById(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	preferentialVO := &vo.PreferentialVO{
		OriginalPrice:   productVO.Price,
		DiscountedPrice: productVO.Price,
		DiscountRate:    1.0,
		DiscountInfo:    "商品金额大于100，触发优惠",
	}

	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			result = preferentialVO
			ctx.Logger().Errorf("recoverErr:%v", recoverErr)
		}
	}()

	if productVO.Price < 100 {
		return preferentialVO, nil
	}

	rule, exists := enum.DiscountRules[userById.WxGrade]

	if !exists {
		ctx.Logger().Infof("不是会员")
		return preferentialVO, nil
	}

	discountedPrice := utils.RoundToTwoDecimalPlaces(productVO.Price * rule.Discount)

	return &vo.PreferentialVO{
		OriginalPrice:     productVO.Price,
		DiscountedPrice:   discountedPrice,
		PreferentialPrice: utils.RoundToTwoDecimalPlaces(productVO.Price - discountedPrice),
		DiscountRate:      rule.Discount,
		DiscountInfo:      fmt.Sprintf("会员等级:%v,折扣:%v折", userById.WxGrade, rule.Discount*100),
	}, nil
}

func (svc OrderService) Finish(ctx jet.Ctx, finishReq *req.OrderFinishReq) error {
	orderPO, _ := svc.orderRepo.FindByID(finishReq.OrderId)
	executorNum := 1
	if orderPO.Executor2Id != -1 {
		executorNum++
	}
	if orderPO.Executor3Id != -1 {
		executorNum++
	}
	// 抽成比例
	cutRate := svc.getCutRate(ctx)
	// 每个人分到的钱
	executorPrice := math.Floor(orderPO.FinalPrice*(1-cutRate)/float64(executorNum)*100) / 100
	err := svc.orderRepo.FinishOrder(ctx, &orderRepoDTO.FinishOrderDTO{
		Id:            finishReq.OrderId,
		Images:        finishReq.Images,
		ExecutorNum:   executorNum,
		ExecutorPrice: executorPrice,
		CutRate:       cutRate,
		OrderInfo:     orderPO,
	})
	if err != nil {
		ctx.Logger().Errorf("[Finish]ERROR: %v", err.Error())
		return errors.New("订单完成失败，请联系客服")
	}
	go func() {
		defer utils.RecoverAndLogError(ctx)
		svc.SendMessagesToExecutors(ctx, orderPO, orderPO.ExecutorID, executorPrice)
		// 给其他打手发送打钱消息
		if orderPO.Executor2Id >= 0 {
			svc.SendMessagesToExecutors(ctx, orderPO, orderPO.Executor2Id, executorPrice)
		}
		if orderPO.Executor3Id >= 0 {
			svc.SendMessagesToExecutors(ctx, orderPO, orderPO.Executor3Id, executorPrice)
		}
		// 检查用户是否需要升级等级了
		svc.userService.checkUserGrade(ctx, orderPO.PurchaseId)
		// 给用户发消息，并提醒其进行评价
		message := fmt.Sprintf(
			"尊敬的老板:您好，您的订单:%v，订单号：%v 已完成，请前往订单列表对打手进行评价",
			orderPO.OrderName,
			orderPO.OrderId,
		)
		_ = svc.messageService.PushSystemMessage(ctx, orderPO.PurchaseId, message)
	}()
	return nil
}

func (svc OrderService) SendMessagesToExecutors(ctx jet.Ctx, orderPO *po.Order, executorID int, executorPrice float64) {
	dashPO, _ := svc.userService.FindUserByDashId(ctx, executorID)
	message := fmt.Sprintf(
		"尊敬的打手:%v(%v)您好，您的订单:%v，订单号：%v 已完成，结算金额：%v",
		dashPO.MemberNumber,
		dashPO.Name,
		orderPO.OrderName,
		orderPO.OrderId,
		executorPrice,
	)
	_ = svc.messageService.PushSystemMessage(ctx, dashPO.ID, message)
}

// getCutRate 返回小数抽成比例
func (svc OrderService) getCutRate(ctx jet.Ctx) (cutRate float64) {
	defer utils.RecoverAndLogError(ctx)
	defer traceUtil.TraceElapsedByName(time.Now(), "getCutRate")
	// 默认抽成20%
	cutRate = 0.2
	cutRatePO, err := svc.commonRepo.FindConfigByName(ctx, commonEnum.CutRate.String())
	if err != nil || cutRatePO == nil {
		ctx.Logger().Errorf("[getCutRate]ERROR: %v", err)
	} else if len(cutRatePO.Content) >= 1 {
		parseString := utils.ParseString(cutRatePO.Content[0]["desc"])
		if utils.IsNumber(parseString) {
			cutRate = utils.ParseFloat64(parseString) / 100
			if cutRate >= 1 {
				cutRate = 0.2
			}
		}
	}
	return
}

func (svc OrderService) GetProcessingOrderList(ctx jet.Ctx) ([]*vo.OrderVO, error) {
	var (
		orders []*po.Order
		err    error
	)
	// 1. 获取金牌打手提前看到订单时间
	dasher, _ := svc.userService.FindUserById(ctx, middleware.MustGetUserId(ctx))
	if dasher.MemberNumber <= 100 {
		orders, err = svc.orderRepo.QueryOrderByStatus(ctx, enum.PROCESSING)
	} else {
		var delayTime int = 20 // 默认20s
		// 非金牌打手 晚指定秒后才能看到订单
		configByName, _ := svc.commonRepo.FindConfigByName(ctx, commonEnum.DelayTime.String())
		if configByName != nil && len(configByName.Content) > 0 && configByName.Content[0]["desc"] != nil {
			delayTime = utils.SafeParseNumber[int](configByName.Content[0]["desc"])
		}
		delayDuration := time.Now().Add(-time.Second * time.Duration(delayTime))
		orders, err = svc.orderRepo.QueryOrderWithDelayTime(ctx, enum.PROCESSING, delayDuration)
	}
	if err != nil {
		ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
		return nil, errors.New("订单查询失败，请联系客服")
	}
	return utils.CopySlice[*po.Order, *vo.OrderVO](orders), nil
}

func (svc OrderService) Start(ctx jet.Ctx, req *req.OrderStartReq) error {
	ctx.Logger().Infof("订单开始:%v", utils.ObjToJsonStr(req))
	err := svc.startOrder(ctx, req.OrderId, req.ExecutorId, req.StartImages)
	if err != nil {
		ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
		return errors.New("订单开始失败，请联系客服")
	}
	go svc.handleLowTimeOutDeduction(ctx, req.OrderId, req.ExecutorId)
	return nil
}

func (svc OrderService) handleLowTimeOutDeduction(ctx jet.Ctx, ordersId uint, executorId int) {
	defer utils.RecoverWithPrefix(ctx, "handleLowTimeOutDeduction")
	// 查询接单时间
	orderPO, err := svc.orderRepo.FindByID(ordersId)
	if err != nil {
		ctx.Logger().Errorf("[FindByID]ERROR: %v", err.Error())
		return
	}
	penaltyStrategy, err := penalty.FetchPenaltyRule(penalty.DeductRuleTimeout)

	if err != nil {
		ctx.Logger().Errorf("fetch penaltyRule ERROR: %v", err)
		return
	}
	applyPenalty, err := penaltyStrategy.ApplyPenalty(
		&penalty.PenaltyReq{
			OrdersId: uint(orderPO.OrderId), GrabTime: orderPO.GrabAt,
		},
	)

	if err != nil || applyPenalty.PenaltyAmount <= 0 {
		ctx.Logger().Errorf("[ApplyPenalty]ERROR: %v", err)
		return
	}
	ctx.Logger().Errorf("deduction applyPenalty sucess: %+v", applyPenalty)

	if applyPenalty.DeductType == penalty.DeductRuleTimeout {

		record, err := svc.deductionRepo.FindByOrderNo(ordersId)

		if err == nil && record != nil && record.ID > 0 {
			// 已经有处罚记录了
			ctx.Logger().Errorf("has durable Deduction info: %+v", utils.ObjToJsonStr(record))
			return
		}

	}

	dasherPO, _ := svc.userService.FindUserByDashId(ctx, executorId)

	err = svc.deductionRepo.InsertOne(&po.Deduction{
		UserID:          dasherPO.ID,
		DasherId:        executorId,
		OrderNo:         ordersId,
		ConfirmPersonId: 0,
		Amount:          applyPenalty.PenaltyAmount,
		Reason:          applyPenalty.Reason,
		Status:          enum.Deduct_PENDING,
	})

	_ = svc.messageService.PushSystemMessage(ctx, dasherPO.ID, applyPenalty.Message)

	if err != nil {
		ctx.Logger().Errorf("deduction insert ERROR: %v", err)
		return
	}
}

func (svc OrderService) startOrder(ctx jet.Ctx, orderId uint, executorId int, image string) error {
	return svc.orderRepo.UpdateOrderByDasher(ctx, orderId, executorId, enum.RUNNING, image)
}

func (svc OrderService) AddOrRemoveExecutor(ctx jet.Ctx, orderReq *req.OrderExecutorReq) (err error) {
	if orderReq.ExecutorName == "" && orderReq.ExecutorId == -1 {
		err = svc.orderRepo.RemoveAssistant(ctx, utils.MustCopy[orderDTO.OrderExecutorDTO](orderReq))
	} else {
		err = svc.orderRepo.AddAssistant(ctx, utils.MustCopy[orderDTO.OrderExecutorDTO](orderReq))
	}
	if err != nil {
		ctx.Logger().Errorf("AddOrRemoveExecutor ERROR:%v", err.Error())
		return errors.New("拒绝失败")
	}
	return
}

// ==================== 提现相关  ====================

func (svc OrderService) HistoryWithDrawAmount(ctx jet.Ctx) (*vo.WithDrawVO, error) {
	userId := middleware.MustGetUserId(ctx)
	userById, err := svc.userService.FindUserById(ctx, userId)
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, cannot find user:%v", userId)
		return nil, errors.New("cannot find user info")
	}
	var (
		approveWithdrawnAmount  float64
		withdrawnAmount         float64
		orderWithdrawAbleAmount float64
		totalDeduct             float64
		wg                      = new(sync.WaitGroup)
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		// 提现成功的钱
		approveWithdrawnAmount, _ = svc.withdrawalRepo.ApproveWithdrawnAmount(ctx, userById.MemberNumber)
		// 四舍五入
		approveWithdrawnAmount = utils.RoundToTwoDecimalPlaces(approveWithdrawnAmount)
	}()
	go func() {
		defer wg.Done()
		// 用户发起提现的钱，包括未提现和提现成功的
		withdrawnAmount, _ = svc.withdrawalRepo.WithdrawnAmountNotReject(ctx, userById.MemberNumber)
		withdrawnAmount = utils.RoundToTwoDecimalPlaces(withdrawnAmount)
	}()
	go func() {
		defer wg.Done()
		// 订单中能提现的钱
		orderWithdrawAbleAmount, _ = svc.orderRepo.OrderWithdrawAbleAmount(ctx, userById.MemberNumber)
		orderWithdrawAbleAmount = utils.RoundToTwoDecimalPlaces(orderWithdrawAbleAmount)
	}()

	go func() {
		defer wg.Done()
		// 罚款的钱
		totalDeduct, _ = svc.deductionRepo.TotalDeduct(ctx, userId)
	}()

	wg.Wait()

	ctx.Logger().Infof(
		"dashId:%v, approveWithdrawnAmount:%v, withdrawnAmount:%v, orderWithdrawAbleAmount:%v,totalDeduct:%v",
		userById.MemberNumber, approveWithdrawnAmount, withdrawnAmount, orderWithdrawAbleAmount, totalDeduct,
	)

	if approveWithdrawnAmount > orderWithdrawAbleAmount {
		ctx.Logger().Errorf(
			"[HistoryWithDrawAmount]ERROR, approveWithdrawnAmount: %v gt orderWithdrawAbleAmount:%v",
			approveWithdrawnAmount, orderWithdrawAbleAmount,
		)
		return nil, errors.New("系统查询错误，请联系管理员")
	}

	minRangeNum, maxRangeNum := svc.fetchWithDrawRange(ctx)

	return &vo.WithDrawVO{
		HistoryWithDrawAmount: utils.RoundToTwoDecimalPlaces(approveWithdrawnAmount),
		WithdrawAbleAmount:    utils.RoundToTwoDecimalPlaces(orderWithdrawAbleAmount - withdrawnAmount - totalDeduct),
		WithdrawRangeMax:      float64(maxRangeNum),
		WithdrawRangeMin:      float64(minRangeNum),
	}, nil
}

func (svc OrderService) fetchWithDrawRange(ctx jet.Ctx) (int, int) {
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

func (svc OrderService) WithDraw(ctx jet.Ctx, drawReq *req.WithDrawReq) error {
	userId := middleware.MustGetUserId(ctx)
	// 1. 添加提现记录
	userById, _ := svc.userService.FindUserById(ctx, userId)
	// 2. 检查提现金额

	err := svc.withdrawalRepo.Withdrawn(ctx, userById.MemberNumber, userId, userById.Name, drawReq.Amount)
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, err:%v", err.Error())
		return errors.New("提现失败，请联系管理员")
	}
	// 2. 发消息，已提交提现申请
	message := fmt.Sprintf("您提现申请已发起，管理员会在3个工作日内处理，您也可以联系管理员进行审批，提现金额为：%v元", drawReq.Amount)
	_ = svc.messageService.PushSystemMessage(ctx, userId, message)
	return nil
}

// GrabOrder 抢单逻辑
func (svc OrderService) GrabOrder(ctx jet.Ctx, grabReq *req.OrderGrabReq) error {
	defer traceUtil.TraceElapsedByName(time.Now(), fmt.Sprintf("%s GrabOrder", ctx.Logger().ReqId))
	// 1. 抢单
	var dasherName string
	dasher, _ := svc.userService.FindUserByDashId(ctx, grabReq.ExecutorId)
	if dasher != nil {
		dasherName = dasher.Name
	}
	err := svc.orderRepo.GrabOrder(ctx, grabReq.OrderId, grabReq.ExecutorId, dasherName)
	if err != nil {
		ctx.Logger().Errorf("[GrabOrder]ERROR, err:%v", err.Error())
		return errors.New("订单已被抢走")
	}
	go func() {
		defer utils.RecoverAndLogError(ctx)
		// 2. 给买家发送消息
		orderPO, _ := svc.orderRepo.FindByID(grabReq.OrderId)
		dasherPO, _ := svc.userService.FindUserByDashId(ctx, grabReq.ExecutorId)
		toUserMessage := fmt.Sprintf(
			"您的订单:%v，已被打手:%v(%v)接受，可前往订单，选中日期%v进行查看",
			orderPO.OrderName, dasherPO.MemberNumber, dasherPO.Name, formatDate(orderPO.PurchaseDate),
		)
		_ = svc.messageService.PushSystemMessage(ctx, orderPO.PurchaseId, toUserMessage)
		// 3. 给打手发消息
		toDasherMessage := fmt.Sprintf(
			"您的订单:%v，已抢单成功，可前往订单，选中日期%v进行查看，请尽快组件队伍开始订单",
			orderPO.OrderName, formatDate(orderPO.PurchaseDate),
		)
		_ = svc.messageService.PushSystemMessage(ctx, dasherPO.ID, toDasherMessage)
	}()
	return nil
}

func formatDate(date *time.Time) string {
	return date.Format("2006-01-02")
}

func (svc OrderService) WithDrawList(ctx jet.Ctx, drawReq *req.WithDrawListReq) ([]*vo.WithDrawListVO, error) {
	query := utils.MustCopy[orderDTO.WithdrawListDTO](drawReq)
	query.UserId = middleware.MustGetUserId(ctx)
	withdrawalRecords, err := svc.withdrawalRepo.ListWithdraw(ctx, query)
	if err != nil {
		ctx.Logger().Errorf("[WithDrawList]ERROR, err:%v", err.Error())
		return nil, errors.New("查询失败")
	}
	vos := utils.CopySlice[*po.WithdrawalRecord, *vo.WithDrawListVO](withdrawalRecords)
	utils.ForEach(vos, func(ele *vo.WithDrawListVO) {
		ele.WithdrawalStatus = enum.WithdrawalStatus(ele.WithdrawalStatus).DisplayName()
	})
	return vos, nil
}

func (svc OrderService) RemoveByID(id int64) error {
	return svc.orderRepo.RemoveByID(id)
}

// ClearAllDasherInfo 清空所有打手信息，重新派单到大厅
func (svc OrderService) ClearAllDasherInfo(ctx jet.Ctx, id uint) error {
	err := svc.orderRepo.ClearOrderDasherInfo(ctx, id)
	if err != nil {
		ctx.Logger().Errorf("[ClearAllDasherInfo]err:%v", err)
		return errors.New("转单失败")
	}
	return nil
}

var (
	syncTimeOutLogger = xlog.NewWith("syncTimeOutLogger")
	dCtx              = xjet.NewDefaultJetContext()
)

// SyncTimeOutOrder 将超时的订单重新发往大厅
func (svc OrderService) SyncTimeOutOrder() {
	// 1. 找到所有打手抢单成功但超时未开始的订单
	orders, err := svc.orderRepo.FindTimeOutOrders(constant.Duration_10_minute)
	if err != nil || orders == nil || len(orders) <= 0 {
		syncTimeOutLogger.Errorf("[SyncTimeOutOrder] ERROR: %v, orders is %+v", err, orders)
		return
	}
	utils.ForEach(orders, func(order *po.Order) {
		defer utils.RecoverAndLogError(dCtx)
		_ = svc.ClearAllDasherInfo(dCtx, order.ID)
		syncTimeOutLogger.Infof("[SyncTimeOutOrder] clear orderInfo, order is: %+v", order)
		// 给打手发送消息
		userPO, _ := svc.userService.FindUserByDashId(dCtx, order.ExecutorID)
		_ = svc.messageService.PushSystemMessage(
			dCtx,
			userPO.ID,
			fmt.Sprintf("您的订单超时未组队，已重新派往接单大厅，订单Id为:%v，如有问题请联系客服", order.OrderId),
		)
	})
}

func (svc OrderService) RemoveAssistantEvent(ctx jet.Ctx) error {
	userId := middleware.MustGetUserId(ctx)
	userPO, _ := svc.userService.FindUserById(ctx, userId)
	return svc.orderRepo.RemoveDasher(ctx, userPO.MemberNumber)
}
