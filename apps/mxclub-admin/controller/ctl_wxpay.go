package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) PostV1WxpayRefunds(ctx jet.Ctx, params *req.WxPayRefundsReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.orderService.Refunds(ctx, params))
}
