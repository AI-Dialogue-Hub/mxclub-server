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
