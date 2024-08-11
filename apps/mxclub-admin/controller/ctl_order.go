package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewOrderController)
}

type OrderController struct {
	jet.BaseJetController
	orderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) jet.ControllerResult {
	return jet.NewJetController(&OrderController{
		orderService: orderService,
	})
}

// =========================================================================

func (ctl OrderController) PostV1OrderList(ctx jet.Ctx, req *req.OrderListReq) (*api.Response, error) {
	pageResult, err := ctl.orderService.List(ctx, req)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) DeleteV1Order0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	orderId, _ := param.GetInt64(0)
	err := ctl.orderService.RemoveByID(orderId)
	return xjet.WrapperResult(ctx, "ok", err)
}

func (ctl OrderController) PostV1Order(ctx jet.Ctx, param *vo.OrderVO) (*api.Response, error) {
	err := ctl.orderService.UpdateOrder(ctx, param)
	return xjet.WrapperResult(ctx, "ok", err)
}

func (ctl OrderController) PostV1OrderTransfer0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	orderId, _ := param.GetInt64(0)
	err := ctl.orderService.TransferOrder(ctx, orderId)
	return xjet.WrapperResult(ctx, "转单成功，订单已回到订单中心", err)
}

// 提现相关

func (ctl OrderController) PostV1WithdrawList(ctx jet.Ctx, req *req.WitchDrawListReq) (*api.Response, error) {
	pageResult, err := ctl.orderService.ListWithdraw(ctx, req)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) PostV1WithdrawUpdate(ctx jet.Ctx, req *req.WitchDrawUpdateReq) (*api.Response, error) {
	err := ctl.orderService.UpdateWithdraw(ctx, req)
	return xjet.WrapperResult(ctx, "ok", err)
}
