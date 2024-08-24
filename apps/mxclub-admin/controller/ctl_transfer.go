package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) GetV1TransferList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	vos, count, err := ctl.orderService.ListTransferInfo(ctx, params)
	pageResult := api.WrapPageResult(params, vos, count)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl OrderController) DeleteV1Transfer0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	got, _ := param.GetInt64(0)
	return xjet.WrapperResult(ctx, "OK", ctl.orderService.RemoveTransfer(ctx, got))
}

func (ctl OrderController) PostV1Transfer(ctx jet.Ctx, vo *vo.TransferVO) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.orderService.UpdateTransfer(ctx, vo))
}
