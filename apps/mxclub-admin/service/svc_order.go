package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewOrderService)
}

type OrderService struct {
	orderRepo     repo.IOrderRepo
	withdrawRepo  repo.IWithdrawalRepo
	deductionRepo repo.IDeductionRepo
}

func NewOrderService(repo repo.IOrderRepo,
	withdrawRepo repo.IWithdrawalRepo,
	deductionRepo repo.IDeductionRepo) *OrderService {
	return &OrderService{orderRepo: repo,
		withdrawRepo:  withdrawRepo,
		deductionRepo: deductionRepo,
	}
}

// =============================================================

func (svc OrderService) List(ctx jet.Ctx, orderReq *req.OrderListReq) (*api.PageResult, error) {
	status := enum.ParseOrderStatusByString(orderReq.OrderStatus)
	list, count, err := svc.orderRepo.ListAroundCache(ctx, orderReq.PageParams, orderReq.Ge, orderReq.Le, status)
	if err != nil {
		ctx.Logger().Errorf("[OrderService]List ERROR:%v", err.Error())
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
		ctx.Logger().Errorf("[OrderService]ListWithdraw ERROR:%v", err.Error())
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
