package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) PostV1DeductionList(ctx jet.Ctx, req *req.DeductionListReq) (*api.Response, error) {
	listDeduction, total, err := ctl.orderService.ListDeduction(ctx, req)
	return xjet.WrapperResult(ctx, api.WrapPageResult(req.PageParams, listDeduction, total), err)
}

func (ctl OrderController) PutV1Deduction(ctx jet.Ctx, req *req.DeductionAddReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, xjet.OK, ctl.orderService.AddDeduction(ctx, req))
}

func (ctl OrderController) PostV1Deduction(ctx jet.Ctx, req *req.DeductionUpdateReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, xjet.OK, ctl.orderService.UpdateDeduction(ctx, req))
}

func (ctl OrderController) DeleteV1Deduction0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	got, _ := param.GetInt64(0)
	return xjet.WrapperResult(ctx, xjet.OK, ctl.orderService.Delete(ctx, uint(got)))
}
