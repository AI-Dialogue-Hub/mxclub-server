package controller

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/domain/common/entity/enum"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewMiniConfigController)
}

type MiniConfigCtr struct {
	jet.BaseJetController
	miniConfigService *service.MiniConfigService
}

func NewMiniConfigController(miniConfigService *service.MiniConfigService) jet.ControllerResult {
	return jet.NewJetController(&MiniConfigCtr{
		miniConfigService: miniConfigService,
	})
}

type putConfigParam struct {
	ConfigName string           `json:"config_name" validate:"required"`
	Content    []map[string]any `json:"content" validate:"required"`
}

func (ctl *MiniConfigCtr) GetV1ConfigList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	list, count, err := ctl.miniConfigService.List(ctx, params)
	return xjet.WrapperResult(ctx, api.WrapPageResult(params, list, count), err)
}

func (ctl *MiniConfigCtr) PutV1Config(ctx jet.Ctx, params *putConfigParam) (*api.Response, error) {
	if enum.MiniConfigEnum(params.ConfigName).IsNotValid() {
		return xjet.WrapperResult(ctx, nil, errors.New("ConfigName is not valid"))
	}
	return xjet.WrapperResult(ctx, "ok", ctl.miniConfigService.AddOrUpdate(ctx, params.ConfigName, params.Content))
}

func (ctl *MiniConfigCtr) DeleteV1Config0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	id := args.CmdArgs[0]
	if xjet.IsAnyEmpty(id) {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "id is empty")
	}
	return xjet.WrapperResult(ctx, "ok", ctl.miniConfigService.Delete(ctx, id))
}

func (ctl *MiniConfigCtr) GetV1Config0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
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
