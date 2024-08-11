package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewOrderService)
}

type OrderService struct {
	orderRepo         repo.IOrderRepo
	wxPayCallbackRepo repo.IWxPayCallbackRepo
	withdrawRepo      repo.IWithdrawalRepo
	deductionRepo     repo.IDeductionRepo
}

func NewOrderService(repo repo.IOrderRepo,
	withdrawRepo repo.IWithdrawalRepo,
	deductionRepo repo.IDeductionRepo,
	wxPayCallbackRepo repo.IWxPayCallbackRepo) *OrderService {
	return &OrderService{orderRepo: repo,
		withdrawRepo:      withdrawRepo,
		deductionRepo:     deductionRepo,
		wxPayCallbackRepo: wxPayCallbackRepo,
	}
}

// =============================================================

func (svc OrderService) List(ctx jet.Ctx, orderReq *req.OrderListReq) (*api.PageResult, error) {
	status := enum.ParseOrderStatusByString(orderReq.OrderStatus)
	list, count, err := svc.orderRepo.ListAroundCache(ctx, orderReq.PageParams, orderReq.Ge, orderReq.Le, status)
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
	if updateReq.WithdrawalStatus == "completed" {

	} else if updateReq.WithdrawalStatus == "reject" {

	}
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", updateReq.Id)
	update.Set("withdrawal_status", updateReq.WithdrawalStatus)
	update.Set("withdrawal_method", updateReq.WithdrawalMethod)
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
	logger.Infof("outRefundNo:%v", outRefundNo)
	err = wxpay.Refunds(ctx, transaction, outRefundNo, params.Reason)
	if err != nil {
		logger.Errorf("err:%v", err)
		return errors.New("退款失败")
	}
	// 3. 修改订单状态
	err = svc.orderRepo.UpdateOrderStatus(ctx, utils.SafeParseUint64(params.OutTradeNo), enum.Refunds)
	if err != nil {
		logger.Errorf("err:%v", err)
		return errors.New("退款失败")
	}
	return nil
}

// TransferOrder 转单就是把
func (svc OrderService) TransferOrder(ctx jet.Ctx, id int64) error {
	err := svc.orderRepo.ClearOrderDasherInfo(ctx, id)
	if err != nil {
		ctx.Logger().Errorf("[TransferOrder]err:%v", err)
		return errors.New("转单失败")
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
