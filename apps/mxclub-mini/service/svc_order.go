package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"math"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
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
		_ = svc.messageService.PushMessage(ctx, dto.NewDispatchMessage(req.ExecutorId, order.ID, req.GameRegion, req.RoleId))
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
			DiscountedPrice: 0,
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
	err := svc.orderRepo.FinishOrder(ctx, finishReq.OrderId, finishReq.Images)
	if err != nil {
		ctx.Logger().Errorf("[Finish]ERROR: %v", err.Error())
		return errors.New("订单完成失败，请联系客服")
	}
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
	// 1. 给其他两个打手发消息
	if req.Executor2Id > 0 {
		_ = svc.messageService.PushMessage(ctx, dto.NewDispatchMessage(req.Executor2Id, req.OrderId, req.GameRegion, req.RoleId))
	}
	if req.Executor3Id > 0 {
		_ = svc.messageService.PushMessage(ctx, dto.NewDispatchMessage(req.Executor3Id, req.OrderId, req.GameRegion, req.RoleId))
	}
	return nil
}

// ==================== 提现相关  ====================

func (svc OrderService) HistoryWithDrawAmount(ctx jet.Ctx) (*vo.WithDrawVO, error) {
	userPO, err := svc.userService.FindUserById(middleware.MustGetUserId(ctx))
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, cannot find user:%v", middleware.MustGetUserId(ctx))
		return nil, errors.New("cannot find user info")
	}
	withdrawnAmount, _ := svc.withdrawalRepo.WithdrawnAmount(ctx, userPO.MemberNumber)
	orderWithdrawAbleAmount, _ := svc.orderRepo.OrderWithdrawAbleAmount(ctx, userPO.MemberNumber)
	if withdrawnAmount > orderWithdrawAbleAmount {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, withdrawnAmount: %v gt orderWithdrawAbleAmount:%v", withdrawnAmount, orderWithdrawAbleAmount)
		return nil, errors.New("系统查询错误，请联系管理员")
	}
	return &vo.WithDrawVO{
		HistoryWithDrawAmount: withdrawnAmount,
		WithdrawAbleAmount:    orderWithdrawAbleAmount - withdrawnAmount,
		WithdrawRangeMax:      20000,
		WithdrawRangeMin:      200,
	}, nil
}

func (svc OrderService) WithDraw(ctx jet.Ctx, drawReq *req.WithDrawReq) error {
	// 1. 发消息，已提交提现申请

	// 2. 添加提现记录
	userPO, _ := svc.userService.FindUserById(middleware.MustGetUserId(ctx))
	err := svc.withdrawalRepo.Withdrawn(ctx, userPO.MemberNumber, drawReq.Amount)
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, err:%v", err.Error())
		return errors.New("提现失败，请联系管理员")
	}
	return nil
}
