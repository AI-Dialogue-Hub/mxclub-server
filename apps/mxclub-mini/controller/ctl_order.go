package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/service"
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

// ============================================================================

func (c *OrderController) PostV1OrderList(ctx jet.Ctx, params *req.OrderListReq) (*api.Response, error) {
	if !params.OrderStatus.Valid() {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "params OrderStatus invalid")
	}
	pageResult, err := c.orderService.List(ctx, params)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (c *OrderController) PostV1WithdrawInfo(ctx jet.Ctx) (*api.Response, error) {
	withDrawVO, err := c.orderService.HistoryWithDrawAmount(ctx)
	return xjet.WrapperResult(ctx, withDrawVO, err)
}

func (c *OrderController) GetV1Preferential0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	productId, _ := param.GetInt64(0)
	preferentialVO, err := c.orderService.Preferential(ctx, uint(productId))
	return xjet.WrapperResult(ctx, preferentialVO, err)
}

func (c *OrderController) PutV1Order(ctx jet.Ctx, req *req.OrderReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", c.orderService.Add(ctx, req))
}
