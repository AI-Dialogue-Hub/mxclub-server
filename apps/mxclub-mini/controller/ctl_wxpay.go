package controller

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewWxPayController)
}

type WxPayController struct {
	jet.BaseJetController
	wxPayService   *service.WxPayService
	messageService *service.MessageService
}

func NewWxPayController(wxPayService *service.WxPayService, messageService *service.MessageService) jet.ControllerResult {
	return jet.NewJetController(&WxPayController{
		wxPayService:   wxPayService,
		messageService: messageService,
	})
}

// ========================================================================================

func (ctl WxPayController) PostWxpayNotify(ctx jet.Ctx, params *maps.LinkedHashMap[string, any]) (*api.Response, error) {
	ctx.Logger().Infof("[PostWxpayNotify] %v", utils.ObjToJsonStr(params))
	go ctl.wxPayService.HandleWxpayNotify(ctx)
	return xjet.WrapperResult(ctx, "ok", nil)
}

func (ctl WxPayController) PostWxpayRefundsNotify(ctx jet.Ctx, params *maps.LinkedHashMap[string, any]) (*api.Response, error) {
	ctx.Logger().Infof("[PostWxpayRefundsNotify] %v", utils.ObjToJsonStr(params))
	return xjet.WrapperResult(ctx, "ok", nil)
}

func (ctl WxPayController) GetWxpayNotify(ctx jet.Ctx, params *maps.LinkedHashMap[string, any]) (*api.Response, error) {
	ctx.Logger().Infof("[PostWxpayNotify] %v", utils.ObjToJsonStr(params))
	go ctl.wxPayService.HandleWxpayNotify(ctx)
	return xjet.WrapperResult(ctx, "ok", nil)
}

func (ctl WxPayController) PostV1WxpayPrepay(ctx jet.Ctx, params *req.WxPayReq) (*api.Response, error) {
	prePayDTO, err := ctl.wxPayService.Prepay(ctx, middleware.MustGetUserId(ctx), params.Amount)
	return xjet.WrapperResult(ctx, prePayDTO, err)
}
