package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) PostV1DeductionList(ctx jet.Ctx, req *req.DeductionListReq) (*api.Response, error) {
	listDeduction, err := ctl.OrderService.ListDeduction(ctx, req)
	return xjet.WrapperResult(ctx, listDeduction, err)
}
