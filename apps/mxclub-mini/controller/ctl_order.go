package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
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

func (ctl OrderController) PostV1OrderList(ctx jet.Ctx, params *req.OrderListReq) (*api.Response, error) {
	if !params.OrderStatus.Valid() {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "params OrderStatus invalid")
	}
	pageResult, err := ctl.orderService.List(ctx, params)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) DeleteV1Order(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	orderId, _ := param.GetInt64(0)
	err := ctl.orderService.RemoveByID(orderId)
	return xjet.WrapperResult(ctx, "ok", err)
}

// PostV1WithdrawInfo 这里是查询可以提现的信息
func (ctl OrderController) PostV1WithdrawInfo(ctx jet.Ctx) (*api.Response, error) {
	withDrawVO, err := ctl.orderService.HistoryWithDrawAmount(ctx)
	return xjet.WrapperResult(ctx, withDrawVO, err)
}

// PostV1Withdraw 进行提现
func (ctl OrderController) PostV1Withdraw(ctx jet.Ctx, drawReq *req.WithDrawReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.WithDraw(ctx, drawReq))
}

// PostV1WithdrawList 提现记录
func (ctl OrderController) PostV1WithdrawList(ctx jet.Ctx, drawReq *req.WithDrawListReq) (*api.Response, error) {
	withDrawList, err := ctl.orderService.WithDrawList(ctx, drawReq)
	return xjet.WrapperResult(ctx, api.WrapPageResult(drawReq.PageParams, withDrawList, 0), err)
}

func (ctl OrderController) GetV1Preferential0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	productId, _ := param.GetInt64(0)
	preferentialVO, err := ctl.orderService.Preferential(ctx, uint(productId))
	return xjet.WrapperResult(ctx, preferentialVO, err)
}

func (ctl OrderController) PutV1Order(ctx jet.Ctx, req *req.OrderReq) (*api.Response, error) {
	ctx.Logger().Infof("user:%v pay success:%v", middleware.MustGetUserId(ctx), utils.ObjToJsonStr(req))
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.Add(ctx, req))
}

func (ctl OrderController) PostV1OrderFinish(ctx jet.Ctx, req *req.OrderFinishReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.Finish(ctx, req))
}

// GetV1OrderDasher 获取抢单大厅里面的订单
func (ctl OrderController) GetV1OrderDasher(ctx jet.Ctx) (*api.Response, error) {
	orderVOS, err := ctl.orderService.GetProcessingOrderList(ctx)
	return xjet.WrapperResult(ctx, orderVOS, err)
}

// PostV1OrderStart 开始订单，如果有其他打手，需要给其他打手发送邀请的消息，如果没有，则进入进行中
func (ctl OrderController) PostV1OrderStart(ctx jet.Ctx, req *req.OrderStartReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.Start(ctx, req))
}

func (ctl OrderController) PostV1OrderExecutorAdd(ctx jet.Ctx, executorReq *req.OrderExecutorReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.AddOrRemoveExecutor(ctx, executorReq))
}

func (ctl OrderController) PostV1OrderExecutorDelete(ctx jet.Ctx, executorReq *req.OrderExecutorReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.AddOrRemoveExecutor(ctx, executorReq))
}

// PostV1OrderGrab 抢单逻辑
func (ctl OrderController) PostV1OrderGrab(ctx jet.Ctx, grabReq *req.OrderGrabReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "抢单成功", ctl.orderService.GrabOrder(ctx, grabReq))
}
