package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewOrderController)
}

type OrderController struct {
	jet.BaseJetController
	OrderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) jet.ControllerResult {
	return jet.NewJetController(&OrderController{
		OrderService: orderService,
	})
}

// =========================================================================

func (ctl OrderController) PostV1OrderList(ctx jet.Ctx, req *req.OrderListReq) (*api.Response, error) {
	pageResult, err := ctl.OrderService.List(ctx, req)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) DeleteV1Order(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	orderId, _ := param.GetInt64(0)
	err := ctl.orderService.orderRepo.RemoveByID(orderId)
	return xjet.WrapperResult(ctx, "ok", err)
}

// 提现相关

func (ctl OrderController) PostV1WithdrawList(ctx jet.Ctx, req *req.WitchDrawListReq) (*api.Response, error) {
	pageResult, err := ctl.OrderService.ListWithdraw(ctx, req)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) PostV1WithdrawUpdate(ctx jet.Ctx, req *req.WitchDrawUpdateReq) (*api.Response, error) {
	err := ctl.OrderService.UpdateWithdraw(ctx, req)
	return xjet.WrapperResult(ctx, "ok", err)
}
