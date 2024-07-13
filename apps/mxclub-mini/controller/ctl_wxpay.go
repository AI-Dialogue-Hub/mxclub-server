package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewWxPayController)
}

type WxPayController struct {
	jet.BaseJetController
	WxPayService   *service.WxPayService
	messageService *service.MessageService
}

func NewWxPayController(wxPayService *service.WxPayService, messageService *service.MessageService) jet.ControllerResult {
	return jet.NewJetController(&WxPayController{
		WxPayService:   wxPayService,
		messageService: messageService,
	})
}

// ========================================================================================

func (ctl WxPayController) PostWxpayNotify(ctx jet.Ctx, params map[string]any) (*api.Response, error) {
	ctx.Logger().Infof("%v", params)
	return xjet.WrapperResult(ctx, "ok", nil)
}

func (ctl WxPayController) GetWxpayNotify(ctx jet.Ctx, params map[string]any) (*api.Response, error) {
	ctx.Logger().Infof("%v", params)
	return xjet.WrapperResult(ctx, "ok", nil)
}
