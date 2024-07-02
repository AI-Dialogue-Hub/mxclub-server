package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/service"
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
	ConfigName string           `json:"config_name"`
	Content    []map[string]any `json:"content"`
}

func (ctr *MiniConfigCtr) GetV1ConfigList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	list, count, err := ctr.miniConfigService.List(ctx, params)
	return xjet.WrapperResult(ctx, api.WrapPageResult(params, list, count), err)
}

func (ctr *MiniConfigCtr) PutV1Config(ctx jet.Ctx, params *putConfigParam) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctr.miniConfigService.Add(ctx, params.ConfigName, params.Content))
}

func (ctr *MiniConfigCtr) DeleteV1Config0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	id := args.CmdArgs[0]
	if xjet.IsAnyEmpty(id) {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "id is empty")
	}
	return xjet.WrapperResult(ctx, "ok", ctr.miniConfigService.Delete(ctx, id))
}
