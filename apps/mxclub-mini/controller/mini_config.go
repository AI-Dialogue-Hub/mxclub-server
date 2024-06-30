package controller

import (
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
)

func init() {
	jet.Provide(NewMiniConfigController)
}

type MiniConfigController struct {
	jet.BaseJetController
	miniConfigService *service.MiniConfigService
}

func NewMiniConfigController(miniConfigService *service.MiniConfigService) jet.ControllerResult {
	return jet.NewJetController(&MiniConfigController{
		miniConfigService: miniConfigService,
	})
}

func (ctl MiniConfigController) GetV1Swiper(ctx jet.Ctx) (*api.Response, error) {
	result, err := ctl.miniConfigService.GetConfigByName(ctx)
	return xjet.WarpperResult(ctx, result, err)
}
