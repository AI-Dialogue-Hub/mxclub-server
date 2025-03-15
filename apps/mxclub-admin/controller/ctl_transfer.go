package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func (ctl OrderController) GetV1TransferList(ctx jet.Ctx, params *req.TransferListReq) (*api.Response, error) {
	vos, count, err := ctl.orderService.ListTransferInfo(ctx, params)
	pageResult := api.WrapPageResult(&params.PageParams, vos, count)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) DeleteV1Transfer0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	got, _ := param.GetInt64(0)
	return xjet.WrapperResult(ctx, "OK", ctl.orderService.RemoveTransfer(ctx, got))
}

func (ctl OrderController) PostV1Transfer(ctx jet.Ctx, vo *vo.TransferVO) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.orderService.UpdateTransfer(ctx, vo))
}

func (ctl OrderController) PostV1TransferTo(ctx jet.Ctx, transferReq *req.TransferReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "转单成功", ctl.orderService.TransferTo(ctx, transferReq))
}

func (ctl OrderController) PostV1OrderTransfer0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	orderId, _ := param.GetInt64(0)
	err := ctl.orderService.ClearAllDasherInfo(ctx, utils.ParseUint(orderId))
	return xjet.WrapperResult(ctx, "转单成功，订单已回到订单中心", err)
}
