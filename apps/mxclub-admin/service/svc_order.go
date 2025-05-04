package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"mxclub/apps/mxclub-admin/entity/dto"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/middleware"
	"mxclub/domain/event"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
	"time"
)

func init() {
	jet.Provide(NewOrderService)
	jet.Invoke(func(u *OrderService) {
		event.RegisterEvent("OrderService", event.EventRemoveDasher, u.RemoveAssistantEvent)
		event.RegisterEvent("TransferService", event.EventRemoveDasher, u.RemoveTransferRecord)
		event.RegisterEvent("DeductService", event.EventRemoveDasher, u.RemoveDeductRecord)
		event.RegisterEvent("WithdrawalService", event.EventRemoveDasher, u.RemoveWithdrawalRecord)
		event.RegisterEvent("EvaluationService", event.EventRemoveDasher, u.RemoveEvaluation)
		event.RegisterEvent("RewardService", event.EventRemoveDasher, u.RemoveRewardRecord)
	})
}

type OrderService struct {
	orderRepo          repo.IOrderRepo
	wxPayCallbackRepo  repo.IWxPayCallbackRepo
	withdrawRepo       repo.IWithdrawalRepo
	deductionRepo      repo.IDeductionRepo
	userRepo           userRepo.IUserRepo
	messageService     *MessageService
	transferRepo       repo.ITransferRepo
	evaluationRepo     repo.IEvaluationRepo
	rewardRepo         repo.IRewardRecordRepo
	operatorLogService *OperatorLogService
}

func NewOrderService(repo repo.IOrderRepo,
	withdrawRepo repo.IWithdrawalRepo,
	deductionRepo repo.IDeductionRepo,
	wxPayCallbackRepo repo.IWxPayCallbackRepo,
	messageService *MessageService,
	userRepo userRepo.IUserRepo,
	transferRepo repo.ITransferRepo,
	evaluationRepo repo.IEvaluationRepo,
	rewardRepo repo.IRewardRecordRepo,
	operatorLogService *OperatorLogService) *OrderService {
	return &OrderService{orderRepo: repo,
		withdrawRepo:       withdrawRepo,
		deductionRepo:      deductionRepo,
		wxPayCallbackRepo:  wxPayCallbackRepo,
		messageService:     messageService,
		userRepo:           userRepo,
		transferRepo:       transferRepo,
		evaluationRepo:     evaluationRepo,
		rewardRepo:         rewardRepo,
		operatorLogService: operatorLogService,
	}
}

// =============================================================

func (svc *OrderService) List(ctx jet.Ctx, orderReq *req.OrderListReq) (*api.PageResult, error) {
	status := enum.ParseOrderStatusByString(orderReq.OrderStatus)
	var (
		orderId    = -1
		executorId = -1
	)
	if orderReq.OrderId != "" {
		orderId = utils.SafeParseNumber[int](orderReq.OrderId)
	}
	if orderReq.ExecutorId > 0 {
		executorId = orderReq.ExecutorId
	}
	// 1. 订单数据
	list, count, err := svc.orderRepo.ListAroundCache(
		ctx, orderReq.PageParams, orderReq.Ge, orderReq.Le, status, orderId, executorId)
	if err != nil {
		ctx.Logger().Errorf("[OrderService]List ERROR:%v", err.Error())
		return nil, errors.New("获取失败")
	}
	orderVOS := utils.CopySlice[*po.Order, *vo.OrderVO](list)
	// 2. 评论数据
	orderIdList := utils.Map[*vo.OrderVO, uint64](orderVOS, func(in *vo.OrderVO) uint64 { return in.OrderId })
	orderId2EvaluationMap, err := svc.evaluationRepo.FindByOrderList(ctx, orderIdList)
	if err == nil && orderId2EvaluationMap != nil && len(orderId2EvaluationMap) > 0 {
		for _, orderVO := range orderVOS {
			if evaluationPOs, ok := orderId2EvaluationMap[orderVO.OrderId]; ok {
				evaluationVOs := utils.CopySlice[*po.OrderEvaluation, *vo.EvaluationVO](evaluationPOs)
				m := make(map[int]*vo.EvaluationVO)
				for _, evaluationVO := range evaluationVOs {
					if evaluationVO.ExecutorID <= 0 && evaluationVO.Rating <= 0 {
						ctx.Logger().Errorf("invalid evaluation => %v", utils.ObjToJsonStr(evaluationVO))
						continue
					}
					m[evaluationVO.ExecutorID] = evaluationVO
				}
				orderVO.EvaluationInfo = m
			}
		}
	}
	// 3. 订单状态
	for _, orderVO := range orderVOS {
		// 3.1 如果状态是已退单，需要拼接退单人信息
		if orderVO.OrderStatus == enum.Refunds {
			refundOrderLog, err := svc.operatorLogService.FindRefundOrderLog(ctx, utils.ParseString(orderVO.OrderId))
			if err != nil || refundOrderLog == nil {
				ctx.Logger().Errorf("[OrderService#List] FindRefundOrderLog ERROR, %v", err)
				// 用默认的兜底
				orderVO.OrderStatusStr = orderVO.OrderStatus.String()
				continue
			}
			orderVO.OrderStatusStr = fmt.Sprintf("%s(%s)", orderVO.OrderStatus.String(), refundOrderLog.Remarks)
		} else {
			orderVO.OrderStatusStr = orderVO.OrderStatus.String()
		}
	}
	return api.WrapPageResult(orderReq.PageParams, orderVOS, count), nil
}

func (svc *OrderService) ListWithdraw(ctx jet.Ctx, params *req.WitchDrawListReq) (*api.PageResult, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(params.Page, params.PageSize)
	if params.WithdrawalStatus != "" && params.WithdrawalStatus != "ALL" {
		query.SetFilter("withdrawal_status = ?", params.WithdrawalStatus)
	}
	if params.DasherId >= 0 {
		query.SetFilter("dasher_id = ?", params.DasherId)
	}
	records, count, err := svc.withdrawRepo.ListByWrapper(ctx, query)
	if err != nil {
		ctx.Logger().Errorf("[*OrderService]ListWithdraw ERROR:%v", err.Error())
		return nil, errors.New("获取失败")
	}
	return api.WrapPageResult(params.PageParams, utils.CopySlice[*po.WithdrawalRecord, *vo.WithdrawVO](records), count), nil
}

func (svc *OrderService) UpdateWithdraw(ctx jet.Ctx, updateReq *req.WitchDrawUpdateReq) error {
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
	update.Set("application_time", time.Now())
	return svc.withdrawRepo.UpdateByWrapper(update)
}

func (svc *OrderService) RemoveByID(ctx jet.Ctx, id int64) error {
	doLogRemoveOperator(ctx, id, svc)
	return svc.orderRepo.RemoveByID(id)
}

func doLogRemoveOperator(ctx jet.Ctx, id int64, svc *OrderService) {
	defer utils.RecoverAndLogError(ctx)
	userInfo, err := middleware.FetchUserInfoByCtx(ctx.FastHttpCtx())
	if err != nil || userInfo == nil {
		ctx.Logger().Errorf("[doLogRemoveOperator] FetchUserInfoByCtx ERROR")
		return
	}
	svc.operatorLogService.DoLog(
		ctx,
		dto.NewOrderRemoveOperatorLogDTO(
			utils.ParseString(id),
			fmt.Sprintf("管理员:%v, 删除订单, 时间:%v", userInfo.Name, time.Now().Format(time.DateTime)),
		),
	)
}

func (svc *OrderService) Refunds(ctx jet.Ctx, params *req.WxPayRefundsReq) error {
	// 0. 日志记录
	doLogRefundsOperatorLog(ctx, params.OutTradeNo, svc)
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
	err = svc.orderRepo.UpdateOrderStatusIncludingDeleted(ctx, utils.SafeParseUint64(params.OutTradeNo), enum.Refunds)
	if err != nil {
		logger.Errorf("[*OrderService#Refunds]err:%v", err)
		return errors.New("退款失败")
	}
	return nil
}

func doLogRefundsOperatorLog(ctx jet.Ctx, orderId string, svc *OrderService) {
	defer utils.RecoverAndLogError(ctx)
	userInfo, err := middleware.FetchUserInfoByCtx(ctx.FastHttpCtx())
	if err != nil || userInfo == nil {
		ctx.Logger().Errorf("[doLogRefundsOperatorLog] FetchUserInfoByCtx ERROR")
		return
	}
	operatorLogDTO := dto.NewOrderRefundsOperatorLogDTO(
		orderId,
		fmt.Sprintf("管理员:%v(%v), 退款订单, 时间:%v", userInfo.Name, userInfo.ID, time.Now().Format(time.DateTime)),
	)
	ctx.Logger().Infof("doLogRefundsOperatorLog: %v", utils.ObjToJsonStr(operatorLogDTO))
	svc.operatorLogService.DoLog(
		ctx,
		operatorLogDTO,
	)
}

func (svc *OrderService) UpdateOrder(ctx jet.Ctx, orderVO *vo.OrderVO) error {
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

func (svc *OrderService) CheckDasherInRunningOrder(ctx jet.Ctx, memberNumber int) bool {
	orderPO, err := svc.orderRepo.FindByDasherId(ctx, memberNumber)
	return err != nil && orderPO != nil && orderPO.ID > 0
}

// ==============  清除打手记录   =====================

func (svc *OrderService) RemoveAssistantEvent(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	return svc.orderRepo.RemoveDasherAllOrderInfo(ctx, userPO.MemberNumber)
}

func (svc *OrderService) RemoveTransferRecord(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	return svc.transferRepo.RemoveByDasherId(ctx, userPO.MemberNumber)
}

func (svc *OrderService) RemoveDeductRecord(ctx jet.Ctx) error {
	return svc.deductionRepo.RemoveDasher(ctx, ctx.MustGet("userId").(uint))
}

func (svc *OrderService) RemoveWithdrawalRecord(ctx jet.Ctx) error {
	return svc.withdrawRepo.RemoveWithdrawalRecord(ctx, ctx.MustGet("userId").(uint))
}

func (svc *OrderService) RemoveEvaluation(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	return svc.evaluationRepo.RemoveEvaluation(ctx, userPO.MemberNumber)
}

// RemoveRewardRecord 清理打手打赏信息
func (svc *OrderService) RemoveRewardRecord(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	if err := svc.rewardRepo.ClearAllRewardByDasherId(ctx, userId); err != nil {
		ctx.Logger().Errorf("RemoveRewardRecord ERROR, %v", err)
		return err
	}
	return nil
}
