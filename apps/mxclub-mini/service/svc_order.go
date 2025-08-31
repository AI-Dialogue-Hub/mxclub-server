package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	traceUtil "github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"math"
	"mxclub/apps/mxclub-mini/config"
	constantMini "mxclub/apps/mxclub-mini/entity/constant"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/event"
	"mxclub/domain/lottery/ability"
	"mxclub/domain/order/biz/penalty"
	orderDTO "mxclub/domain/order/entity/dto"
	productRepo "mxclub/domain/product/repo"
	userPOInfo "mxclub/domain/user/po"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/common/txsms"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/constant"
	"runtime/debug"
	"strconv"
	"strings"
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

// setUp 初始化注册handle
//
// @see MessageService
// @see RewardService
func setUp(svc *OrderService) {
	logger := xlog.NewWith("RegisterEvent")
	logger.Info("RegisterEvent =====>>>> start")
	event.RegisterEvent("OrderService", event.EventRemoveDasher, svc.RemoveAssistantEvent)
	event.RegisterEvent("TransferService", event.EventRemoveDasher, svc.RemoveTransferRecord)
	event.RegisterEvent("DeductService", event.EventRemoveDasher, svc.RemoveDeductRecord)
	event.RegisterEvent("WithdrawalService", event.EventRemoveDasher, svc.RemoveWithdrawalRecord)
	event.RegisterEvent("EvaluationService", event.EventRemoveDasher, svc.RemoveEvaluation)
	logger.Info("RegisterEvent =====>>>> end")
}

type OrderService struct {
	orderRepo            repo.IOrderRepo
	withdrawalRepo       repo.IWithdrawalRepo
	userService          *UserService
	productService       *ProductService
	messageService       *MessageService
	commonRepo           commonRepo.IMiniConfigRepo
	deductionRepo        repo.IDeductionRepo
	evaluationRepo       repo.IEvaluationRepo
	transferRepo         repo.ITransferRepo
	productSalesRepo     productRepo.IProductSalesRepo
	rewardRecordRepo     repo.IRewardRecordRepo
	wxNotifyService      *WxNotifyService
	lotteryAbility       ability.ILotteryAbility
	deactivateDasherRepo userRepo.IDeactivateDasherRepo
	wxPayCallbackRepo    repo.IWxPayCallbackRepo
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
	transferRepo repo.ITransferRepo,
	productSalesRepo productRepo.IProductSalesRepo,
	rewardRecordRepo repo.IRewardRecordRepo,
	wxNotifyService *WxNotifyService,
	lotteryAbility ability.ILotteryAbility,
	deactivateDasherRepo userRepo.IDeactivateDasherRepo,
	wxPayCallbackRepo repo.IWxPayCallbackRepo) *OrderService {
	orderService = &OrderService{
		orderRepo:            repo,
		withdrawalRepo:       withdrawalRepo,
		userService:          userService,
		productService:       productService,
		messageService:       messageService,
		commonRepo:           commonRepo,
		deductionRepo:        deductionRepo,
		evaluationRepo:       evaluationRepo,
		transferRepo:         transferRepo,
		productSalesRepo:     productSalesRepo,
		rewardRecordRepo:     rewardRecordRepo,
		wxNotifyService:      wxNotifyService,
		lotteryAbility:       lotteryAbility,
		deactivateDasherRepo: deactivateDasherRepo,
		wxPayCallbackRepo:    wxPayCallbackRepo,
	}
	// 初始化
	setUp(orderService)
	return orderService
}

// ===============================================================

func (svc *OrderService) Add(ctx jet.Ctx, req *req.OrderReq) error {
	if _, err := svc.AddByOrderStatus(ctx, req, enum.PrePay); err != nil {
		return errors.New("订单已经添加成功")
	}
	return nil
}

// PaySuccessOrder 下单成功
func (svc *OrderService) PaySuccessOrder(ctx jet.Ctx, orderNo uint64) error {
	defer utils.RecoverAndLogError(ctx)
	if orderNo <= 0 {
		ctx.Logger().Errorf("invalid orderNo: %v", orderNo)
		return errors.New("invalid orderNo")
	}
	err := svc.orderRepo.UpdateOrderStatusIncludingDeleted(ctx, orderNo, enum.PROCESSING)
	if err != nil {
		ctx.Logger().Errorf("UpdateOrderStatus failed, err: %v", err)
		return errors.New("更新失败")
	}
	orderPO, _ := svc.orderRepo.FindByOrderOrOrdersId(ctx, uint(orderNo))
	var (
		executorId   int
		dasherPO     *userPOInfo.User
		orderTradeNo = utils.SafeParseUint64(orderPO.OrderId)
	)
	// 增加销量
	if err = svc.productSalesRepo.AddOrUpdateSale(ctx, orderPO.ProductID, productRepo.Default_Sale_Volume); err != nil {
		ctx.Logger().Errorf("[PaySuccessOrder]sales add failed, err:%v", err)
	}
	// 如果指定打手，给打手发送接单消息
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
				ctx.Logger().Infof(
					"指定订单, id:%v, dasherId:%v,order: %v", orderPO.OrderId, dasherPO.MemberNumber, orderPO)
				// 发送派单信息
				svc.messageService.PushMessage(
					ctx,
					dto.NewDispatchMessageWithFinalPrice(
						dasherPO.ID, uint(orderTradeNo), orderPO.GameRegion, orderPO.RoleId, "",
						utils.RoundToTwoDecimalPlaces(orderPO.FinalPrice)),
				)
				// 微信消息通知
				//_ = svc.wxNotifyService.SendMessage(ctx, dasherPO.ID, "您有新的指定订单，请赶快前往小程序查看!")
				// 短信推送消息
				err = txsms.SendDefaultDispatchMsg(dasherPO.Phone)
				if err != nil {
					ctx.Logger().Errorf("[PaySuccessOrder] phone:%v send sms failed:%v", dasherPO.Phone, err)
				}
			}()
		}
	} else {
		executorId = -1
	}
	ctx.Logger().Infof("pay success, order: %v", utils.ObjToJsonStr(orderPO))
	return nil
}

func (svc *OrderService) AddByOrderStatus(ctx jet.Ctx, req *req.OrderReq, status enum.OrderStatus) (*po.Order, error) {
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
	// 更新电话 和 游戏Id
	if req.Phone != "" {
		go func() {
			_ = svc.userService.userRepo.UpdateUserPhone(ctx, userId, req.Phone)
			if id, extractErr := po.ExtractID(req.RoleId); extractErr == nil {
				_ = svc.userService.userRepo.UpdateUserGameId(ctx, userId, id)
			}
		}()
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

func (svc *OrderService) List(ctx jet.Ctx, req *req.OrderListReq) (*api.PageResult, error) {
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
	// 将gameId和roleId分开
	id2OrderPOMap := utils.SliceToSingleMap(list, func(ele *po.Order) uint { return ele.ID })
	for _, orderVO := range orderVOS {
		if orderPO, ok := id2OrderPOMap.Get(orderVO.ID); ok {
			orderVO.RoleId = orderPO.FetchRoleId()
			orderVO.GameId = orderPO.FetchGameId()
		}
	}

	// 用户可以打赏的状态
	orderIds := utils.Map(orderVOS, func(in *vo.OrderVO) string {
		// 先全部置于可以打赏的状态
		in.CanReward = true
		return utils.ParseString(in.OrderId)
	})
	orderId2RewardsMap, err := svc.rewardRecordRepo.FindByOrderIds(ctx, orderIds)
	if err == nil && orderId2RewardsMap != nil && len(orderId2RewardsMap) > 0 {
		for _, orderVO := range orderVOS {
			records, ok := orderId2RewardsMap[utils.ParseString(orderVO.OrderId)]
			if ok && records != nil && len(records) > 0 {
				orderVO.CanReward = false
			} else {
				orderVO.CanReward = true
			}
		}
	}
	return api.WrapPageResult(&req.PageParams, orderVOS, 0), err
}

func (svc *OrderService) doBuildUserGrade(ctx jet.Ctx, vos []*vo.OrderVO) {
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

func (svc *OrderService) Preferential(ctx jet.Ctx, productId uint) (result *vo.PreferentialVO, err error) {
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

func (svc *OrderService) Finish(ctx jet.Ctx, finishReq *req.OrderFinishReq) error {
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

func (svc *OrderService) SendMessagesToExecutors(ctx jet.Ctx, orderPO *po.Order, executorID int, executorPrice float64) {
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
func (svc *OrderService) getCutRate(ctx jet.Ctx) (cutRate float64) {
	return getCutRate(ctx, svc.commonRepo)
}

func getCutRate(ctx jet.Ctx, commonRepo commonRepo.IMiniConfigRepo) (cutRate float64) {
	defer utils.RecoverAndLogError(ctx)
	defer traceUtil.TraceElapsedByName(time.Now(), "getCutRate")
	// 默认抽成20%
	cutRate = 0.2
	cutRatePO, err := commonRepo.FindConfigByName(ctx, commonEnum.CutRate.String())
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

func (svc *OrderService) GetProcessingOrderList(ctx jet.Ctx) ([]*vo.OrderVO, error) {
	var (
		orders           []*po.Order
		err              error
		defaultDelayTime = 20 // 默认delay 20s
	)
	fetchDelayTime := func(index int) (delayDate time.Time) {
		if r := recover(); r != nil {
			ctx.Logger().Error("Recovered from panic:", r)
			debug.PrintStack()
			return time.Now().Add(-time.Second * time.Duration(defaultDelayTime))
		}
		var delayTime int
		configByName, _ := svc.commonRepo.FindConfigByName(ctx, commonEnum.DelayTime.String())
		if configByName != nil && len(configByName.Content) >= index && configByName.Content[index]["desc"] != nil {
			delayTime = utils.SafeParseNumber[int](configByName.Content[0]["desc"])
		} else {
			ctx.Logger().Infof("fetch delayTime error, set default value %v", defaultDelayTime)
			delayTime = defaultDelayTime
		}
		delayDuration := time.Now().Add(-time.Second * time.Duration(delayTime))
		return delayDuration
	}
	// 1. 获取金牌打手提前看到订单时间
	dasher, _ := svc.userService.FindUserById(ctx, middleware.MustGetUserId(ctx))
	payName := config.GetConfig().WxPayConfig.PayName
	// 金牌打手可以及时看到所有订单 三角洲金牌逻辑和其他不一样，其他小程序是打手id小于100的
	if (payName == "明星三角洲" && dasher.DasherLevel == userEnum.DasherLevel_Gold) ||
		(payName != "明星三角洲" && dasher.MemberNumber < 100) {
		orders, err = svc.orderRepo.QueryOrderByStatus(ctx, enum.PROCESSING)
	} else if dasher.DasherLevel == userEnum.DasherLevel_Silver {
		// 银牌打手延迟多久看到订单
		silverDasherDelayTime := fetchDelayTime(0)
		ctx.Logger().Infof("silverDasherDelayTime:%v", silverDasherDelayTime)
		orders, err = svc.orderRepo.QueryOrderWithDelayTime(ctx, enum.PROCESSING, silverDasherDelayTime)
	} else if dasher.DasherLevel == userEnum.DasherLevel_Bronze {
		// 铜牌打手延迟多久看到订单
		dasherDelayTime := fetchDelayTime(1)
		ctx.Logger().Infof("silverDasherDelayTime:%v", dasherDelayTime)
		orders, err = svc.orderRepo.QueryOrderWithDelayTime(ctx, enum.PROCESSING, dasherDelayTime)
	} else {
		// 其他打手
		dasherDelayTime := fetchDelayTime(2)
		ctx.Logger().Infof("silverDasherDelayTime:%v", dasherDelayTime)
		orders, err = svc.orderRepo.QueryOrderWithDelayTime(ctx, enum.PROCESSING, dasherDelayTime)
	}
	if err != nil {
		ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
		return nil, errors.New("订单查询失败，请联系客服")
	}
	orderVOS := utils.CopySlice[*po.Order, *vo.OrderVO](orders)
	// 转盘单要隐藏名称
	utils.ForEach(orderVOS, func(ele *vo.OrderVO) {
		if ability.IsLotteryOrder(strconv.FormatUint(ele.OrderId, 10)) {
			ele.OrderName = fmt.Sprintf("%s，接单后显示内容", extractPrefixBeforeColon(ele.OrderName))
		}
	})
	// 将gameId和roleId分开
	id2OrderPOMap := utils.SliceToSingleMap(orders, func(ele *po.Order) uint { return ele.ID })
	for _, orderVO := range orderVOS {
		if orderPO, ok := id2OrderPOMap.Get(orderVO.ID); ok {
			orderVO.RoleId = orderPO.FetchRoleId()
			orderVO.GameId = orderPO.FetchGameId()
		}
	}
	return orderVOS, nil
}

func extractPrefixBeforeColon(input string) string {
	// 查找冒号的位置
	colonIndex := strings.Index(input, ":")
	if colonIndex == -1 {
		return input // 如果没有冒号，返回整个字符串
	}

	// 返回冒号前的部分
	return input[:colonIndex]
}

func (svc *OrderService) Start(ctx jet.Ctx, req *req.OrderStartReq) error {
	ctx.Logger().Infof("订单开始:%v", utils.ObjToJsonStr(req))
	// 0. 需要确保订单不会被打两次
	orderPO, err := svc.orderRepo.FindByOrderOrOrdersId(ctx, req.OrderId)
	if err != nil {
		ctx.Logger().Errorf(
			"[OrderService#Start]cannot find order, "+
				"req:%v, invalid order:%v", utils.ObjToJsonStr(req), utils.ObjToJsonStr(orderPO))
		return errors.New("订单不存在")
	}
	if orderPO.CompletionDate != nil {
		ctx.Logger().Errorf(
			"[OrderService#Start]req:%v, invalid order:%v", utils.ObjToJsonStr(req), utils.ObjToJsonStr(orderPO))
		return errors.New("该订单状态异常，请联系管理员")
	}
	// 1. 确保订单还在当前打手手上
	if orderPO.ExecutorID != req.ExecutorId {
		ctx.Logger().Errorf(
			"[OrderService#Start]current dasher is not first executor, req:%v, order:%v",
			utils.ObjToJsonStr(req), utils.ObjToJsonStr(orderPO))
		return errors.New("您不是该订单的车头，请检查订单状态")
	}

	err = svc.startOrder(ctx, req.OrderId, req.ExecutorId, req.StartImages)

	if err != nil {
		ctx.Logger().Errorf("[GetProcessingOrderList]ERROR: %v", err.Error())
		return errors.New("订单开始失败，请联系客服")
	}
	go svc.handleLowTimeOutDeduction(ctx, req.OrderId, req.ExecutorId)
	return nil
}

func (svc *OrderService) handleLowTimeOutDeduction(ctx jet.Ctx, ordersId uint, executorId int) {
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
			OrdersId: uint(orderPO.OrderId), OrderRawPrice: orderPO.OriginalPrice, GrabTime: orderPO.GrabAt,
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

func (svc *OrderService) startOrder(ctx jet.Ctx, orderId uint, executorId int, image string) error {
	return svc.orderRepo.UpdateOrderByDasher(ctx, orderId, executorId, enum.RUNNING, image)
}

func (svc *OrderService) AddOrRemoveExecutor(ctx jet.Ctx, orderReq *req.OrderExecutorReq) (err error) {
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

func (svc *OrderService) HistoryWithDrawAmount(ctx jet.Ctx) (*vo.WithDrawVO, error) {
	userId := middleware.MustGetUserId(ctx)
	userById, err := svc.userService.FindUserById(ctx, userId)
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, cannot find user:%v", userId)
		return nil, errors.New("cannot find user info")
	}
	if userById.Role != userEnum.RoleAssistant || userById.MemberNumber < 0 {
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
		approveWithdrawnAmount, _ = svc.withdrawalRepo.ApproveWithdrawnAmount(ctx, userById.MemberNumber)
		// 四舍五入
		approveWithdrawnAmount = utils.RoundToTwoDecimalPlaces(approveWithdrawnAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		// 用户发起提现的钱，包括未提现和提现成功的
		withdrawnAmount, _ = svc.withdrawalRepo.WithdrawnAmountNotReject(ctx, userById.MemberNumber)
		withdrawnAmount = utils.RoundToTwoDecimalPlaces(withdrawnAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		// 订单中能提现的钱
		orderWithdrawAbleAmount, _ = svc.orderRepo.OrderWithdrawAbleAmount(ctx, userById.MemberNumber)
		orderWithdrawAbleAmount = utils.RoundToTwoDecimalPlaces(orderWithdrawAbleAmount)
	}()

	go func() {
		defer utils.RecoverAndLogError(ctx)
		defer wg.Done()
		rewardAmount, _ = svc.rewardRecordRepo.AllRewardAmountByDasherId(ctx, userById.ID)
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
		userById.MemberNumber, approveWithdrawnAmount,
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

func (svc *OrderService) WithDraw(ctx jet.Ctx, drawReq *req.WithDrawReq) error {
	userId := middleware.MustGetUserId(ctx)
	userById, _ := svc.userService.FindUserById(ctx, userId)
	// 0. 查找是否有未完成的提现记录 如果有则限制此次提现
	records, err := svc.withdrawalRepo.FindWithdrawnByStatus(ctx, userById.MemberNumber, enum.Initiated())
	if records != nil && len(records) > 0 {
		ctx.Logger().Errorf("[WithDraw] has durable withdraw records => %v", utils.ObjToJsonStr(records))
		return errors.New("您还有其他提现记录未完结，请结束后再进行提现")
	}
	// 1. 检查提现金额
	minAmount, _ := svc.fetchWithDrawRange(ctx)
	if drawReq.Amount < float64(minAmount) {
		ctx.Logger().Errorf("withDraw Amount:%v less more minAmount:%v", drawReq.Amount, minAmount)
		return errors.New(fmt.Sprintf("提现金额不能小于最小限制:%v", minAmount))
	}
	// 2. 检查余额是否充足
	withDrawVO, err := svc.HistoryWithDrawAmount(ctx)
	if err != nil {
		return err
	}
	if withDrawVO != nil && withDrawVO.WithdrawAbleAmount < drawReq.Amount {
		ctx.Logger().Errorf(
			"withdrawAbleAmount Amount:%v less more withdrawAmount:%v",
			drawReq.Amount, withDrawVO.WithdrawAbleAmount)
		return errors.New("账户余额不足")
	}
	err = svc.withdrawalRepo.Withdrawn(ctx, userById.MemberNumber, userId, userById.Name, drawReq.Amount)
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, err:%v", err.Error())
		return errors.New("提现失败，请联系管理员")
	}
	// 3. 发消息，已提交提现申请
	message := fmt.Sprintf("您提现申请已发起，管理员会在3个工作日内处理，"+
		"您也可以联系管理员进行审批，提现金额为：%v元，请带图联系董事长现结，勿催", drawReq.Amount)
	_ = svc.messageService.PushSystemMessage(ctx, userId, message)
	return nil
}

// GrabOrder 抢单逻辑
func (svc *OrderService) GrabOrder(ctx jet.Ctx, grabReq *req.OrderGrabReq) error {
	defer traceUtil.TraceElapsedByName(time.Now(), fmt.Sprintf("%s GrabOrder", ctx.Logger().ReqId))
	// 0. 最多有两笔进行中订单 且十分钟内只能抢一单
	orders, _ := svc.orderRepo.FindByDasherIdAndStatus(ctx, grabReq.ExecutorId, enum.PROCESSING, enum.RUNNING)
	if orders != nil && len(orders) > 1 {
		ctx.Logger().Errorf("[GrabOrder] has durable orders => %v", utils.ObjToJsonStr(orders))
		return errors.New("进行中订单不能超过两单")
	}
	if orders != nil && len(orders) == 1 {
		orderPO := orders[0]
		if time.Duration(time.Now().Unix()-orderPO.GrabAt.Unix()).Minutes() < 10 {
			ctx.Logger().Errorf(
				"[GrabOrder] has durable orders between 10 minutes => %v", utils.ObjToJsonStr(orders))
			return errors.New("您有进行中的订单，十分钟内只能抢一单")
		}
	}
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
		"您的订单:%v，订单号:%v,已抢单成功，可前往订单，选中日期%v进行查看，请尽快组件队伍开始订单",
		orderPO.OrderName, orderPO.OrderId, formatDate(orderPO.PurchaseDate),
	)
	_ = svc.messageService.PushSystemMessage(ctx, dasherPO.ID, toDasherMessage)
	ctx.Logger().Infof("GrabOrder success, dasherId:%v, orderId:%v", dasher.MemberNumber, orderPO.OrderId)
	return nil
}

func formatDate(date *time.Time) string {
	return date.Format("2006-01-02")
}

func (svc *OrderService) WithDrawList(ctx jet.Ctx, drawReq *req.WithDrawListReq) ([]*vo.WithDrawListVO, error) {
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

func (svc *OrderService) RemoveByID(id int64) error {
	return svc.orderRepo.RemoveByID(id)
}

// ClearAllDasherInfo 清空所有打手信息，重新派单到大厅
func (svc *OrderService) ClearAllDasherInfo(ctx jet.Ctx, id uint) error {
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
func (svc *OrderService) SyncTimeOutOrder() {
	// 1. 找到所有打手抢单成功但超时未开始的订单
	orders, err := svc.orderRepo.FindTimeOutOrders(constant.Duration_10_minute)
	if err != nil || orders == nil || len(orders) <= 0 {
		syncTimeOutLogger.Errorf("[SyncTimeOutOrder] ERROR: %v, orders is %+v", err, orders)
		return
	}
	utils.ForEach(orders, func(order *po.Order) {
		defer utils.RecoverAndLogError(dCtx)
		if order.CompletionDate != nil {
			syncTimeOutLogger.Errorf(
				"[SyncTimeOutOrder#]has order CompletionDate is not nil:%v", utils.ObjToJsonStr(order))
			return
		}
		_ = svc.ClearAllDasherInfo(dCtx, order.ID)
		syncTimeOutLogger.Infof("[SyncTimeOutOrder] clear orderInfo, order is: %+v", order)
		// 给打手发送消息
		userPO, err := svc.userService.FindUserByDashId(dCtx, order.ExecutorID)
		if err != nil || userPO == nil || userPO.ID <= 0 {
			syncTimeOutLogger.Errorf("[SyncTimeOutOrder] FindUserByDashId ERROR, %v, userPO:%v", err, userPO)
			return
		}
		_ = svc.messageService.PushSystemMessage(
			dCtx,
			userPO.ID,
			fmt.Sprintf("您的订单超时未组队，已重新派往接单大厅，订单Id为:%v，如有问题请联系客服", order.OrderId),
		)
		// 清理打手接单的派单消息
		if order.SpecifyExecutor && order.ExecutorID >= 0 {
			var userByDashId *userPOInfo.User
			if userByDashId, err =
				svc.userService.FindUserByDashId(dCtx, order.ExecutorID); err == nil && userByDashId != nil {
				// 清除指定打手的接单消息
				if err = svc.messageService.ClearDispatchMessage(order.OrderId, userByDashId.ID); err != nil {
					syncTimeOutLogger.Errorf("[SyncTimeOutOrder#ClearDispatchMessage]error:%v", err)
				} else {
					syncTimeOutLogger.Infof("[SyncTimeOutOrder#ClearDispatchMessage]SUCCESS, orderInfo:%v",
						utils.ObjToJsonStr(order))
				}
			}
		}
		// 清理指定打手的派单消息
		if order.OutRefundNo != "" {
			// 清除指定打手的接单消息
			err = svc.messageService.ClearDispatchMessage(order.OrderId, utils.SafeParseNumber[uint](order.OutRefundNo))
			if err != nil {
				syncTimeOutLogger.Errorf("[SyncTimeOutOrder#ClearDispatchMessage]OutRefundNo message error:%v", err)
			} else {
				syncTimeOutLogger.Infof("[SyncTimeOutOrder#ClearDispatchMessage]OutRefundNo message "+
					"SUCCESS, orderInfo:%v",
					utils.ObjToJsonStr(order))
			}
		}
	})
}

func (svc *OrderService) RemoveAssistantEvent(ctx jet.Ctx) error {
	if value, exists := ctx.Get(constantMini.LOGOUT_DASHER_ID); exists {
		return svc.orderRepo.RemoveDasherAllOrderInfo(ctx, value.(int))
	} else {
		userId := middleware.MustGetUserId(ctx)
		userPO, _ := svc.userService.FindUserById(ctx, userId)
		// 0. 注销前，打印账户余额信息
		if historyWithDrawAmount, err := svc.HistoryWithDrawAmount(ctx); err == nil {
			ctx.Logger().Infof("[RemoveAssistantEvent] dasher:%v, info:%v, HistoryWithDrawAmount info => %v",
				userPO.MemberNumber, utils.ObjToJsonStr(userPO), utils.ObjToJsonStr(historyWithDrawAmount))
			go func() {
				defer utils.RecoverAndLogError(ctx)
				ctx.Logger().Infof("[RemoveAssistantEvent] dasher:%v, info:%v, HistoryWithDrawAmount info => %v",
					userPO.MemberNumber, utils.ObjToJsonStr(userPO), utils.ObjToJsonStr(historyWithDrawAmount))
				allOrderPOList, _ := svc.orderRepo.FindAllByDasherId(ctx, userPO.MemberNumber)
				// 保存打手最后的金额
				_ = svc.deactivateDasherRepo.InsertOne(&userPOInfo.DeactivateDasher{
					DasherID:              userPO.MemberNumber,
					DasherName:            userPO.Name,
					HistoryWithdrawAmount: historyWithDrawAmount.HistoryWithDrawAmount,
					WithdrawAbleAmount:    historyWithDrawAmount.WithdrawAbleAmount,
					OrderSnapshot:         utils.ObjToJsonStr(allOrderPOList),
				})
			}()
		}
		return svc.orderRepo.RemoveDasherAllOrderInfo(ctx, userPO.MemberNumber)
	}
}

func (svc *OrderService) Refunds(ctx jet.Ctx, params *req.WxPayRefundsReq) (string, error) {
	// 0. 日志记录
	doLogRefundsOperatorLog(ctx, params.OrderId)
	var (
		logger     = ctx.Logger()
		successMsg = "退款请求已发起，请等待客服进行处理"
	)
	// 0.1 检查订单状态
	orderPO, err := svc.orderRepo.FindByOrderOrOrdersId(ctx, utils.ParseUint(params.OrderId))
	if orderPO.OrderStatus != enum.PROCESSING || orderPO.ExecutorID > 0 {
		ctx.Logger().Infof("[OrderService#Refunds]order status error, orderId:%v", params.OrderId)
		return "", errors.New("进行中订单无法退款，可联系客服退款")

	}
	orderId := utils.ParseString(orderPO.OrderId)
	// 1. 查询回调的参数
	wxPayCallbackInfo, err := svc.wxPayCallbackRepo.FindByTraceNo(orderId)
	if err != nil {
		logger.Errorf("err:%v", err)
		return "", errors.New("查询不到对应订单信息")
	}
	// 1.2 转换
	transaction := utils.MustMapToObj[payments.Transaction](wxPayCallbackInfo.RawData)
	// 2. 进行退款
	outRefundNo := wxpay.GenerateOutRefundNo()
	logger.Infof("outRefundNo:%v, orderId:%v", outRefundNo, orderId)
	if params.Reason == "" {
		params.Reason = "协商一致退款"
	}
	err = wxpay.Refunds(ctx, transaction, outRefundNo, params.Reason)
	if err != nil {
		logger.Errorf("err:%v", err)
		return "", errors.New("退款失败")
	}
	// 3. 修改订单状态
	err = svc.orderRepo.UpdateOrderStatusIncludingDeleted(ctx, utils.SafeParseUint64(orderId), enum.Refunds)
	if err != nil {
		logger.Errorf("[*OrderService#Refunds]err:%v", err)
		return "", errors.New("退款失败")
	}
	return successMsg, nil
}

func doLogRefundsOperatorLog(ctx jet.Ctx, orderId string) {
	defer utils.RecoverAndLogError(ctx)
	userId := middleware.MustGetUserId(ctx)
	ctx.Logger().Infof("doLogRefundsOperatorLog: %v, orderId:%v", userId, orderId)

}
