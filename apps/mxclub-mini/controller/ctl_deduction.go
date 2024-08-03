package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) PostV1DeductionList(ctx jet.Ctx, req *req.DeductionListReq) (*api.Response, error) {
	listDeduction, err := ctl.orderService.ListDeduction(ctx, req)
	return xjet.WrapperResult(ctx, api.WrapPageResult(req.PageParams, listDeduction, 0), err)
}
