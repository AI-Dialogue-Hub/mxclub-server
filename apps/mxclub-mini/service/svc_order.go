package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewOrderService)
}

type OrderService struct {
	orderRepo repo.IOrderRepo
}

func NewOrderService(repo repo.IOrderRepo) *OrderService {
	return &OrderService{orderRepo: repo}
}

func (svc OrderService) Add(ctx jet.Ctx, req *req.OrderReq) error {
	return nil
}

func (svc OrderService) List(ctx jet.Ctx, req *req.OrderListReq) (*api.PageResult, error) {
	list, count, err := svc.orderRepo.ListByOrderStatus(ctx, req.OrderStatus, &req.PageParams)
	return api.WrapPageResult(&req.PageParams, utils.CopySlice[*po.Order, *vo.OrderVO](list), count), err
}
