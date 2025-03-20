package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewWxNotifyController)
}

func NewWxNotifyController(notifyService *service.WxNotifyService) jet.ControllerResult {
	return jet.NewJetController(&WxNotifyController{
		notifyService: notifyService,
	})
}

// WxNotifyController 微信消息推送
type WxNotifyController struct {
	jet.BaseJetController
	notifyService *service.WxNotifyService
}

// ======================================================================================

func (ctl WxNotifyController) GetWxMessageNotifyStatus(ctx jet.Ctx, templateReq *req.NotifyTemplateReq) (*api.Response, error) {
	status := ctl.notifyService.FindSubStatus(ctx, templateReq.TemplateId)
	return xjet.WrapperResult(ctx, status, nil)
}

func (ctl WxNotifyController) PutWxMessageNotify(ctx jet.Ctx, templateReq *req.NotifyTemplateReq) (*api.Response, error) {
	err := ctl.notifyService.AddSubNotifyRecord(ctx, templateReq.TemplateId)
	return xjet.WrapperResult(ctx, "ok", err)
}

func (ctl WxNotifyController) PostWxMessageNotifySend(ctx jet.Ctx, req *req.WxNotifySendReq) (*api.Response, error) {
	err := ctl.notifyService.SendMessage(ctx, req.UserId, req.Message)
	return xjet.WrapperResult(ctx, "ok", err)
}
