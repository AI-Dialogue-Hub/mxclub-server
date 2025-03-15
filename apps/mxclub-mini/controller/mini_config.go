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

func (ctr MiniConfigController) GetV1Config0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	if args == nil || xjet.IsAnyEmpty(args.CmdArgs...) {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "configName is empty")
	}
	configName := args.CmdArgs[0]
	if enum.MiniConfigEnum(configName).IsNotValid() {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "config type is not valid")
	}
	result, err := ctr.miniConfigService.GetConfigByName(ctx, configName)
	return xjet.WrapperResult(ctx, result, err)
}

func (ctr MiniConfigController) GetV1ProductSellingpoint(ctx jet.Ctx) (*api.Response, error) {
	return xjet.WrapperResult(ctx, ctr.miniConfigService.FetchSellingPoints(ctx), nil)
}

func (ctr MiniConfigController) GetV1ConfigList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	list, count, err := ctr.miniConfigService.List(ctx, params)
	return xjet.WrapperResult(ctx, api.WrapPageResult(params, list, count), err)
}
