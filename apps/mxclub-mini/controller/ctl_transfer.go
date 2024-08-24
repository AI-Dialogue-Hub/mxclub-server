package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) PostV1Transfer(ctx jet.Ctx, req *req.TransferReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, xjet.OK, ctl.orderService.Transfer(ctx, req))
}
