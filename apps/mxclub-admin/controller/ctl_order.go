package controller

import (
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
)

func init() {
	jet.Provide(NewOrderController)
}

type OrderController struct {
	jet.BaseJetController
	orderService *service.OrderService
	excelService *service.ExcelService
}

func NewOrderController(orderService *service.OrderService, excelService *service.ExcelService) jet.ControllerResult {
	return jet.NewJetController(&OrderController{
		orderService: orderService,
		excelService: excelService,
	})
}

// =========================================================================

func (ctl OrderController) PostV1OrderList(ctx jet.Ctx, req *req.OrderListReq) (*api.Response, error) {
	pageResult, err := ctl.orderService.List(ctx, req)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) DeleteV1Order0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	orderId, _ := param.GetInt64(0)
	err := ctl.orderService.RemoveByID(ctx, orderId)
	return xjet.WrapperResult(ctx, "ok", err)
}

func (ctl OrderController) PostV1Order(ctx jet.Ctx, param *vo.OrderVO) (*api.Response, error) {
	err := ctl.orderService.UpdateOrder(ctx, param)
	return xjet.WrapperResult(ctx, "ok", err)
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

func (ctl OrderController) PostV1ExportTrade(ctx jet.Ctx, req *req.OrderTradeExportReq) {
	ctl.excelService.ExportExcel(ctx, req.StartDate, req.EndDate)
}

func (ctl OrderController) PostV1OrderHistorywithdrawamount(ctx jet.Ctx, req *req.HistoryWithDrawAmountReq) (*api.Response, error) {
	historyWithDrawAmount, err := ctl.orderService.HistoryWithDrawAmount(ctx, req)
	return xjet.WrapperResult(ctx, historyWithDrawAmount, err)
}

func (ctl OrderController) PostV1OrderWithdrawamountAll(ctx jet.Ctx) {
	_ = ctl.orderService.ExportAllDasherWithDrawAmount(ctx)
}

func (ctl OrderController) PostV1OrderDeactivateDasherList(ctx jet.Ctx, req *req.DeactivateReq) (*api.Response, error) {
	pageResult, count, err := ctl.orderService.ListDeactivateDasher(ctx, req)
	return xjet.WrapperResult(ctx, api.WrapPageResult(req.PageParams, pageResult, count), err)
}
