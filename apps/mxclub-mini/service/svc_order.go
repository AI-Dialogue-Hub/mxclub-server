package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/product/po"
	"mxclub/domain/product/repo"
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
	list, err := svc.orderRepo.ListByOrderStatus(ctx, req.OrderStatus, &req.PageParams, req.Ge, req.Le)
	return api.WrapPageResult(&req.PageParams, utils.CopySlice[*po.Order, *vo.OrderVO](list), 0), err
}
