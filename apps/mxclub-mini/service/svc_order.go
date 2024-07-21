package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	traceUtil "github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"math"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	orderDTO "mxclub/domain/order/entity/dto"

	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/wxpay"
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
}

func NewOrderService(
	repo repo.IOrderRepo,
	withdrawalRepo repo.IWithdrawalRepo,
	userService *UserService,
	productService *ProductService,
	messageService *MessageService) *OrderService {

	return &OrderService{
		orderRepo:      repo,
		withdrawalRepo: withdrawalRepo,
		userService:    userService,
		productService: productService,
		messageService: messageService,
	}
}

// ===============================================================

func (svc OrderService) Add(ctx jet.Ctx, req *req.OrderReq) error {
	userId := middleware.MustGetUserId(ctx)
	_ = svc.userService.userRepo.UpdateUserPhone(ctx, userId, req.Phone)
	// 1. 查询商品信息
	// 1.1 折扣信息
	preferentialVO, _ := svc.Preferential(ctx, req.ProductId)
	// 2. 创建订单
	order := &po.Order{
		OrderId:         utils.ParseUint64(wxpay.GenerateUniqueOrderNumber()),
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
		ExecutorID:      req.ExecutorId,
		Notes:           req.Notes,
		DiscountPrice:   preferentialVO.OriginalPrice - preferentialVO.DiscountedPrice,
		FinalPrice:      preferentialVO.DiscountedPrice,
		ExecutorPrice:   0,
		PurchaseDate:    core.Time(time.Now()),
	}
	// 3. 保存订单
	err := svc.orderRepo.InsertOne(order)
	if err != nil {
		ctx.Logger().Errorf("[OrderService]Add ERROR, %v", err.Error())
		ctx.Logger().Errorf("order:%v", utils.ObjToJsonStr(order))
		return errors.New("订单保存保存失败，请联系客服")
	}
	// 4. 如果指定订单，给打手发送接单消息
	if req.SpecifyExecutor {
		_ = svc.messageService.PushMessage(ctx, dto.NewDispatchMessage(req.ExecutorId, order.ID, req.GameRegion, req.RoleId, ""))
	}
	return nil
}

func (svc OrderService) List(ctx jet.Ctx, req *req.OrderListReq) (*api.PageResult, error) {
	userId := middleware.MustGetUserId(ctx)
	list, err := svc.orderRepo.ListByOrderStatus(ctx, req.OrderStatus, &req.PageParams, req.Ge, req.Le, req.MemberNumber, userId)
	if err != nil {
		ctx.Logger().Errorf("[OrderService]List ERROR, %v", err.Error())
		return nil, errors.New("查询不到数据")
	}
	orderVOS := utils.CopySlice[*po.Order, *vo.OrderVO](list)
	return api.WrapPageResult(&req.PageParams, orderVOS, 0), err
}

func (svc OrderService) Preferential(ctx jet.Ctx, productId uint) (*vo.PreferentialVO, error) {
	userId := middleware.MustGetUserId(ctx)
	userPO, err := svc.userService.FindUserById(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	productVO, err := svc.productService.FindById(productId)
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

	rule, exists := enum.DiscountRules[userPO.WxGrade]

	if !exists {
		return nil, errors.New("不是会员")
	}

	discountedPrice := math.Floor(productVO.Price*rule.Discount*100) / 100

	return &vo.PreferentialVO{
		OriginalPrice:     productVO.Price,
		DiscountedPrice:   discountedPrice,
		PreferentialPrice: productVO.Price - discountedPrice,
		DiscountRate:      rule.Discount,
		DiscountInfo:      fmt.Sprintf("会员等级:%v,折扣:%v折", userPO.WxGrade, rule.Discount*100),
	}, nil
}

func (svc OrderService) Finish(ctx jet.Ctx, finishReq *req.OrderFinishReq) error {
	orderPO, _ := svc.orderRepo.FindByID(finishReq.OrderId)
	executorNum := 1
	if orderPO.Executor2Id != 0 {
		executorNum++
	}
	if orderPO.Executor3Id != 0 {
		executorNum++
	}
	// 每个人分到的钱
	executorPrice := math.Floor(orderPO.FinalPrice*0.8/float64(executorNum)*100) / 100
	err := svc.orderRepo.FinishOrder(ctx, finishReq.OrderId, finishReq.Images, executorNum, executorPrice)
	if err != nil {
		ctx.Logger().Errorf("[Finish]ERROR: %v", err.Error())
		return errors.New("订单完成失败，请联系客服")
	}
	go func() {
		defer utils.RecoverAndLogError(ctx)
		// 车头
		dashPO, _ := svc.userService.FindUserByDashId(orderPO.ExecutorID)
		message := fmt.Sprintf(
			"尊敬的打手:%v(%v)您好，您的订单:%v，订单号：%v 已完成，结算金额：%v",
			dashPO.MemberNumber,
			dashPO.Name,
			orderPO.OrderName,
			orderPO.OrderId,
			executorPrice,
		)
		_ = svc.messageService.PushSystemMessage(ctx, dashPO.ID, message)
		// 给其他打手发送打钱消息
		if orderPO.Executor2Id != 0 {
			dash2PO, _ := svc.userService.FindUserByDashId(orderPO.Executor2Id)
			message = fmt.Sprintf(
				"尊敬的打手:%v(%v)您好，您的订单:%v，订单号：%v 已完成，结算金额：%v",
				dash2PO.MemberNumber,
				dash2PO.Name,
				orderPO.OrderName,
				orderPO.OrderId,
				executorPrice,
			)
			_ = svc.messageService.PushSystemMessage(ctx, dash2PO.ID, message)
		}
		if orderPO.Executor3Id != 0 {
			dash3PO, _ := svc.userService.FindUserByDashId(orderPO.Executor3Id)
			message = fmt.Sprintf(
				"尊敬的打手:%v(%v)您好，您的订单:%v，订单号：%v 已完成，结算金额：%v",
				dash3PO.MemberNumber,
				dash3PO.Name,
				orderPO.OrderName,
				orderPO.OrderId,
				executorPrice,
			)
			_ = svc.messageService.PushSystemMessage(ctx, dash3PO.ID, message)
		}
	}()
	return nil
}

func (svc OrderService) GetProcessingOrderList(ctx jet.Ctx) ([]*vo.OrderVO, error) {
	orders, err := svc.orderRepo.QueryOrderByStatus(ctx, enum.PROCESSING)
	if err != nil {
		ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
		return nil, errors.New("订单完成失败，请联系客服")
	}
	return utils.CopySlice[*po.Order, *vo.OrderVO](orders), nil
}

func (svc OrderService) Start(ctx jet.Ctx, req *req.OrderStartReq) error {
	if req.Executor2Id == 0 && req.Executor3Id == 0 {
		// 直接开始
		err := svc.startOrder(ctx, req.OrderId, req.ExecutorId)
		if err != nil {
			ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
			return errors.New("订单开始失败，请联系客服")
		}
		return nil
	}
	// 指定打手的数量
	executorNumber := 0
	if req.Executor2Id > 0 {
		executorNumber++
	}
	if req.Executor3Id > 0 {
		executorNumber++
	}
	// 1. 给其他两个打手发消息
	if req.Executor2Id > 0 {
		user1, _ := svc.userService.FindUserByDashId(req.Executor2Id)
		message := dto.NewDispatchMessage(user1.ID, req.OrderId, req.GameRegion, req.RoleId, utils.ParseString(executorNumber))
		_ = svc.messageService.PushMessage(ctx, message)
	}
	if req.Executor3Id > 0 {
		user2, _ := svc.userService.FindUserByDashId(req.Executor3Id)
		message := dto.NewDispatchMessage(user2.ID, req.OrderId, req.GameRegion, req.RoleId, utils.ParseString(executorNumber))
		_ = svc.messageService.PushMessage(ctx, message)
	}
	return nil
}

func (svc OrderService) startOrder(ctx jet.Ctx, orderId uint, executorId uint) error {
	return svc.orderRepo.UpdateOrderByDasher(ctx, orderId, executorId, enum.RUNNING)
}

func (svc OrderService) AddOrRemoveExecutor(ctx jet.Ctx, orderReq *req.OrderExecutorReq) (err error) {
	if orderReq.ExecutorName == "" && orderReq.ExecutorId == 0 {
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
	userPO, err := svc.userService.FindUserById(middleware.MustGetUserId(ctx))
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, cannot find user:%v", middleware.MustGetUserId(ctx))
		return nil, errors.New("cannot find user info")
	}
	var (
		approveWithdrawnAmount  float64
		withdrawnAmount         float64
		orderWithdrawAbleAmount float64
		c1                      = make(chan struct{})
		c2                      = make(chan struct{})
		c3                      = make(chan struct{})
	)
	go func() {
		defer func() { c1 <- struct{}{} }()
		approveWithdrawnAmount, _ = svc.withdrawalRepo.ApproveWithdrawnAmount(ctx, userPO.MemberNumber)
	}()
	go func() {
		defer func() { c2 <- struct{}{} }()
		withdrawnAmount, _ = svc.withdrawalRepo.WithdrawnAmountNotReject(ctx, userPO.MemberNumber)
	}()
	go func() {
		defer func() { c3 <- struct{}{} }()
		orderWithdrawAbleAmount, _ = svc.orderRepo.OrderWithdrawAbleAmount(ctx, userPO.MemberNumber)
	}()

	<-c1
	<-c2
	<-c3

	if approveWithdrawnAmount > orderWithdrawAbleAmount {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, approveWithdrawnAmount: %v gt orderWithdrawAbleAmount:%v", approveWithdrawnAmount, orderWithdrawAbleAmount)
		return nil, errors.New("系统查询错误，请联系管理员")
	}
	return &vo.WithDrawVO{
		HistoryWithDrawAmount: approveWithdrawnAmount,
		WithdrawAbleAmount:    orderWithdrawAbleAmount - withdrawnAmount,
		WithdrawRangeMax:      20000,
		WithdrawRangeMin:      200,
	}, nil
}

func (svc OrderService) WithDraw(ctx jet.Ctx, drawReq *req.WithDrawReq) error {
	userId := middleware.MustGetUserId(ctx)
	// 1. 添加提现记录
	userPO, _ := svc.userService.FindUserById(userId)
	err := svc.withdrawalRepo.Withdrawn(ctx, userPO.MemberNumber, drawReq.Amount)
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
		return errors.New("抢单失败，请刷新订单列表")
	}
	go func() {
		defer utils.RecoverAndLogError(ctx)
		// 2. 给买家发送消息
		orderPO, _ := svc.orderRepo.FindByID(grabReq.OrderId)
		dasherPO, _ := svc.userService.FindUserByDashId(grabReq.ExecutorId)
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
