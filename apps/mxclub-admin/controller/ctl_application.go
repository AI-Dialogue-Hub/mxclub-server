package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewApplicationController)
}

type ApplicationController struct {
	jet.BaseJetController
	applicationService *service.ApplicationService
}

func NewApplicationController(app *service.ApplicationService) jet.ControllerResult {
	return jet.NewJetController(&ApplicationController{
		applicationService: app,
	})
}

func (ctl ApplicationController) PostV1ApplicationList(ctx jet.Ctx, params *req.ApplicationListReq) (*api.Response, error) {
	pageResult, err := ctl.applicationService.List(ctx, params)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl ApplicationController) PostV1Application(ctx jet.Ctx, req *req.ApplicationReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "Ok", ctl.applicationService.UpdateStatus(ctx, req))
}
