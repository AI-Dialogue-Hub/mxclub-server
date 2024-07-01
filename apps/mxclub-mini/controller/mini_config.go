package controller

import (
	"mxclub/apps/mxclub-mini/service"
	"mxclub/domain/common/entity/enum"
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

func (ctl MiniConfigController) GetV1Config0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	if xjet.IsNil(args) || xjet.IsAnyEmpty(args.CmdArgs...) {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "configName is empty")
	}
	configName := args.CmdArgs[0]
	if enum.MiniConfigEnum(configName).IsNotValid() {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "config type is not valid")
	}
	result, err := ctl.miniConfigService.GetConfigByName(ctx, configName)
	return xjet.WrapperResult(ctx, result, err)
}
