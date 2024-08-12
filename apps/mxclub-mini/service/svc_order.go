package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	traceUtil "github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"math"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	orderDTO "mxclub/domain/order/entity/dto"
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
}

func NewOrderService(
	repo repo.IOrderRepo,
	withdrawalRepo repo.IWithdrawalRepo,
	userService *UserService,
	productService *ProductService,
	messageService *MessageService,
	commonRepo commonRepo.IMiniConfigRepo,
	deductionRepo repo.IDeductionRepo,
	evaluationRepo repo.IEvaluationRepo) *OrderService {

	return &OrderService{
		orderRepo:      repo,
		withdrawalRepo: withdrawalRepo,
		userService:    userService,
		productService: productService,
		messageService: messageService,
		commonRepo:     commonRepo,
		deductionRepo:  deductionRepo,
		evaluationRepo: evaluationRepo,
	}
}

// ===============================================================

func (svc OrderService) Add(ctx jet.Ctx, req *req.OrderReq) error {
	userId := middleware.MustGetUserId(ctx)
	go func() { _ = svc.userService.userRepo.UpdateUserPhone(ctx, userId, req.Phone) }()
	var (
		dasherName            string
		executorId            int
		specifyExecutorUserId uint
	)
	if req.SpecifyExecutor {
		// 特殊编号打手
		executorId = req.ExecutorId
		if executorId == -1 {
			executorId = 0
		}
		dasher, _ := svc.userService.FindUserByDashId(ctx, executorId)
		dasherName = dasher.Name
		specifyExecutorUserId = dasher.ID
	}
	// 1.2 折扣信息
	preferentialVO, _ := svc.Preferential(ctx, req.ProductId)
	// 2. 创建订单
	order := &po.Order{
		OrderId:         utils.SafeParseUint64(req.OrderTradeNo),
		PurchaseId:      userId,
		OrderName:       req.OrderName,
		OrderIcon:       req.OrderIcon,
		OrderStatus:     enum.PROCESSING,
		OriginalPrice:   preferentialVO.OriginalPrice,
		ProductID:       req.ProductId,
		Phone:           req.Phone,
		GameRegion:      req.GameRegion,
		RoleId:          req.RoleId,
		SpecifyExecutor: req.SpecifyExecutor,
		ExecutorID:      executorId,
		Executor2Id:     -1,
		Executor3Id:     -1,
		ExecutorName:    dasherName,
		Notes:           req.Notes,
		DiscountPrice:   preferentialVO.OriginalPrice - preferentialVO.DiscountedPrice,
		FinalPrice:      preferentialVO.DiscountedPrice,
		ExecutorPrice:   0,
		PurchaseDate:    utils.Ptr(time.Now()),
	}
	// 3. 保存订单
	err := svc.orderRepo.InsertOne(order)
	if err != nil {
		ctx.Logger().Errorf("[orderService]Add ERROR, %v", err.Error())
		ctx.Logger().Errorf("order:%v", utils.ObjToJsonStr(order))
		return errors.New("订单保存保存失败，请联系客服")
	}
	// 4. 如果指定订单，给打手发送接单消息
	if req.SpecifyExecutor {
		go svc.messageService.PushMessage(ctx, dto.NewDispatchMessage(specifyExecutorUserId, order.ID, req.GameRegion, req.RoleId, ""))
	}
	return nil
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

func (svc OrderService) Preferential(ctx jet.Ctx, productId uint) (*vo.PreferentialVO, error) {
	userId := middleware.MustGetUserId(ctx)
	userById, err := svc.userService.FindUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	productVO, err := svc.productService.FindById(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	if productVO.Price < 100 {
		return &vo.PreferentialVO{
			OriginalPrice:   productVO.Price,
			DiscountedPrice: productVO.Price,
			DiscountRate:    1.0,
			DiscountInfo:    "商品金额大于100，不触发优惠",
		}, nil
	}

	rule, exists := enum.DiscountRules[userById.WxGrade]

	if !exists {
		return nil, errors.New("不是会员")
	}

	discountedPrice := math.Floor(productVO.Price*rule.Discount*100) / 100

	return &vo.PreferentialVO{
		OriginalPrice:     productVO.Price,
		DiscountedPrice:   discountedPrice,
		PreferentialPrice: productVO.Price - discountedPrice,
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
		_ = svc.messageService.PushSystemMessage(ctx, orderPO.ProductID, message)
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
		return nil, errors.New("订单完成失败，请联系客服")
	}
	return utils.CopySlice[*po.Order, *vo.OrderVO](orders), nil
}

func (svc OrderService) Start(ctx jet.Ctx, req *req.OrderStartReq) error {
	ctx.Logger().Infof("订单开始:%v", utils.ObjToJsonStr(req))
	err := svc.startOrder(ctx, req.OrderId, req.ExecutorId)
	if err != nil {
		ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
		return errors.New("订单开始失败，请联系客服")
	}
	return nil
}

func (svc OrderService) startOrder(ctx jet.Ctx, orderId uint, executorId int) error {
	return svc.orderRepo.UpdateOrderByDasher(ctx, orderId, executorId, enum.RUNNING)
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
		approveWithdrawnAmount, _ = svc.withdrawalRepo.ApproveWithdrawnAmount(ctx, userById.MemberNumber)
	}()
	go func() {
		defer wg.Done()
		withdrawnAmount, _ = svc.withdrawalRepo.WithdrawnAmountNotReject(ctx, userById.MemberNumber)
	}()
	go func() {
		defer wg.Done()
		orderWithdrawAbleAmount, _ = svc.orderRepo.OrderWithdrawAbleAmount(ctx, userById.MemberNumber)
	}()

	go func() {
		defer wg.Done()
		totalDeduct, _ = svc.deductionRepo.TotalDeduct(ctx, userId)
	}()

	wg.Wait()

	ctx.Logger().Infof(
		"dashId:%v, approveWithdrawnAmount:%v, withdrawnAmount:%v, orderWithdrawAbleAmount:%v,totalDeduct:%v",
		userById.MemberNumber, approveWithdrawnAmount, withdrawnAmount, orderWithdrawAbleAmount, totalDeduct,
	)

	if approveWithdrawnAmount > orderWithdrawAbleAmount {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, approveWithdrawnAmount: %v gt orderWithdrawAbleAmount:%v", approveWithdrawnAmount, orderWithdrawAbleAmount)
		return nil, errors.New("系统查询错误，请联系管理员")
	}
	return &vo.WithDrawVO{
		HistoryWithDrawAmount: approveWithdrawnAmount,
		WithdrawAbleAmount:    orderWithdrawAbleAmount - withdrawnAmount - totalDeduct,
		WithdrawRangeMax:      20000,
		WithdrawRangeMin:      200,
	}, nil
}

func (svc OrderService) WithDraw(ctx jet.Ctx, drawReq *req.WithDrawReq) error {
	userId := middleware.MustGetUserId(ctx)
	// 1. 添加提现记录
	userById, _ := svc.userService.FindUserById(ctx, userId)
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
	err := svc.orderRepo.GrabOrder(ctx, grabReq.OrderId, grabReq.ExecutorId)
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
