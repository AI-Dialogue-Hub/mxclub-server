package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/event"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewOrderService)
	jet.Invoke(func(u *OrderService) {
		event.RegisterEvent("OrderService", event.EventRemoveDasher, u.RemoveAssistantEvent)
		event.RegisterEvent("TransferService", event.EventRemoveDasher, u.RemoveTransferRecord)
		event.RegisterEvent("DeductService", event.EventRemoveDasher, u.RemoveDeductRecord)
		event.RegisterEvent("WithdrawalService", event.EventRemoveDasher, u.RemoveWithdrawalRecord)
		event.RegisterEvent("EvaluationService", event.EventRemoveDasher, u.RemoveEvaluation)
	})
}

type OrderService struct {
	orderRepo         repo.IOrderRepo
	wxPayCallbackRepo repo.IWxPayCallbackRepo
	withdrawRepo      repo.IWithdrawalRepo
	deductionRepo     repo.IDeductionRepo
	userRepo          userRepo.IUserRepo
	messageService    *MessageService
	transferRepo      repo.ITransferRepo
	evaluationRepo    repo.IEvaluationRepo
}

func NewOrderService(repo repo.IOrderRepo,
	withdrawRepo repo.IWithdrawalRepo,
	deductionRepo repo.IDeductionRepo,
	wxPayCallbackRepo repo.IWxPayCallbackRepo,
	messageService *MessageService,
	userRepo userRepo.IUserRepo,
	transferRepo repo.ITransferRepo,
	evaluationRepo repo.IEvaluationRepo) *OrderService {
	return &OrderService{orderRepo: repo,
		withdrawRepo:      withdrawRepo,
		deductionRepo:     deductionRepo,
		wxPayCallbackRepo: wxPayCallbackRepo,
		messageService:    messageService,
		userRepo:          userRepo,
		transferRepo:      transferRepo,
		evaluationRepo:    evaluationRepo,
	}
}

// =============================================================

func (svc OrderService) List(ctx jet.Ctx, orderReq *req.OrderListReq) (*api.PageResult, error) {
	status := enum.ParseOrderStatusByString(orderReq.OrderStatus)
	var orderId = -1
	if orderReq.OrderId != "" {
		orderId = utils.SafeParseNumber[int](orderReq.OrderId)
	}
	list, count, err := svc.orderRepo.ListAroundCache(
		ctx, orderReq.PageParams, orderReq.Ge, orderReq.Le, status, orderId)
	if err != nil {
		ctx.Logger().Errorf("[orderService]List ERROR:%v", err.Error())
		return nil, errors.New("获取失败")
	}
	orderVOS := utils.CopySlice[*po.Order, *vo.OrderVO](list)
	utils.ForEach(orderVOS, func(vo *vo.OrderVO) {
		vo.OrderStatusStr = vo.OrderStatus.String()
	})
	return api.WrapPageResult(orderReq.PageParams, orderVOS, count), nil
}

func (svc OrderService) ListWithdraw(ctx jet.Ctx, params *req.WitchDrawListReq) (*api.PageResult, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(params.Page, params.PageSize)
	if params.WithdrawalStatus != "" && params.WithdrawalStatus != "ALL" {
		query.SetFilter("withdrawal_status = ?", params.WithdrawalStatus)
	}
	records, count, err := svc.withdrawRepo.ListByWrapper(ctx, query)
	if err != nil {
		ctx.Logger().Errorf("[orderService]ListWithdraw ERROR:%v", err.Error())
		return nil, errors.New("获取失败")
	}
	return api.WrapPageResult(params.PageParams, utils.CopySlice[*po.WithdrawalRecord, *vo.WithdrawVO](records), count), nil
}

func (svc OrderService) UpdateWithdraw(ctx jet.Ctx, updateReq *req.WitchDrawUpdateReq) error {
	dasherPO, err := svc.userRepo.FindByMemberNumber(ctx, updateReq.DasherId)
	if err != nil {
		return errors.New("打手id错误")
	}
	if updateReq.WithdrawalStatus == "completed" {
		// 发送提现成功消息
		message := fmt.Sprintf("你的提现已通过，提现金额: %v 元，打款方式为：%v", updateReq.WithdrawalAmount, updateReq.WithdrawalInfo)
		_ = svc.messageService.PushSystemMessage(ctx, dasherPO.ID, message)
	} else if updateReq.WithdrawalStatus == "reject" {
		message := fmt.Sprintf("你的提现申请被拒绝，请联系客服或重新发起提现，提现金额: %v 元，拒绝原因：%v",
			updateReq.WithdrawalAmount, updateReq.WithdrawalInfo)
		_ = svc.messageService.PushSystemMessage(ctx, dasherPO.ID, message)
	}
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", updateReq.Id)
	update.Set("withdrawal_status", updateReq.WithdrawalStatus)
	update.Set("withdrawal_method", updateReq.WithdrawalInfo)
	return svc.withdrawRepo.UpdateByWrapper(update)
}

func (svc OrderService) RemoveByID(id int64) error {
	return svc.orderRepo.RemoveByID(id)
}

func (svc OrderService) Refunds(ctx jet.Ctx, params *req.WxPayRefundsReq) error {
	logger := ctx.Logger()
	// 1. 查询回调的参数
	wxPayCallbackInfo, err := svc.wxPayCallbackRepo.FindByTraceNo(params.OutTradeNo)
	if err != nil {
		logger.Errorf("err:%v", err)
		return errors.New("查询不到对应订单信息")
	}
	// 1.2 转换
	transaction := utils.MustMapToObj[payments.Transaction](wxPayCallbackInfo.RawData)
	// 2. 进行退款
	outRefundNo := wxpay.GenerateOutRefundNo()
	logger.Infof("outRefundNo:%v, orderId:%v", outRefundNo, params.OutTradeNo)
	err = wxpay.Refunds(ctx, transaction, outRefundNo, params.Reason)
	if err != nil {
		logger.Errorf("err:%v", err)
		return errors.New("退款失败")
	}
	// 3. 修改订单状态
	err = svc.orderRepo.UpdateOrderStatus(ctx, utils.SafeParseUint64(params.OutTradeNo), enum.Refunds)
	if err != nil {
		logger.Errorf("[OrderService#Refunds]err:%v", err)
		return errors.New("退款失败")
	}
	return nil
}

func (svc OrderService) UpdateOrder(ctx jet.Ctx, orderVO *vo.OrderVO) error {
	updateMap := utils.ObjToMap(orderVO)
	delete(updateMap, "id")
	delete(updateMap, "order_status_str")
	delete(updateMap, "detail_images")
	err := svc.orderRepo.UpdateById(updateMap, orderVO.ID)
	if err != nil {
		ctx.Logger().Errorf("[UpdateOrder]err:%v", err)
		return errors.New("更新订单失败")
	}
	svc.orderRepo.ClearOrderCache(ctx)
	return nil
}

func (svc OrderService) CheckDasherInRunningOrder(ctx jet.Ctx, memberNumber int) bool {
	orderPO, err := svc.orderRepo.FindByDasherId(ctx, memberNumber)
	return err != nil && orderPO != nil && orderPO.ID > 0
}

// ==============  清除打手记录   =====================

func (svc OrderService) RemoveAssistantEvent(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	return svc.orderRepo.RemoveDasher(ctx, userPO.MemberNumber)
}

func (svc OrderService) RemoveTransferRecord(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	return svc.transferRepo.RemoveByDasherId(ctx, userPO.MemberNumber)
}

func (svc OrderService) RemoveDeductRecord(ctx jet.Ctx) error {
	return svc.deductionRepo.RemoveDasher(ctx, ctx.MustGet("userId").(uint))
}

func (svc OrderService) RemoveWithdrawalRecord(ctx jet.Ctx) error {
	return svc.withdrawRepo.RemoveWithdrawalRecord(ctx, ctx.MustGet("userId").(uint))
}

func (svc OrderService) RemoveEvaluation(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	return svc.evaluationRepo.RemoveEvaluation(ctx, userPO.MemberNumber)
}
