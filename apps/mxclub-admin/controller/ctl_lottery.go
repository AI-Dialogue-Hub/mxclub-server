package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewLotteryController)
}

type LotteryController struct {
	jet.BaseJetController
	LotteryService *service.LotteryService
}

func NewLotteryController(app *service.LotteryService) jet.ControllerResult {
	return jet.NewJetController(&LotteryController{LotteryService: app})
}

func (ctl *LotteryController) GetV1LotteryPrizeType(ctx jet.Ctx) (*api.R[*vo.LotteryTypeVO], error) {
	return api.Ok(ctx.Logger().ReqId, ctl.LotteryService.FetchLotteryPrizeType()), nil
}

func (ctl *LotteryController) PostV1LotteryPrizeList(ctx jet.Ctx, params *req.LotteryPrizePageReq) (*api.Response, error) {
	listPrize, count, err := ctl.LotteryService.ListPrize(ctx, params)
	pageResult := api.WrapPageResult(params.PageParams, listPrize, count)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl *LotteryController) PostV1LotteryPrize(ctx jet.Ctx, req *req.LotteryPrizeReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.LotteryService.AddOrUpdatePrize(ctx, req))
}

func (ctl *LotteryController) PostV1LotteryPrizeDel(ctx jet.Ctx, req *req.LotteryPrizeReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.LotteryService.DelPrize(ctx, req))
}

// =============================  ability  =======================================

func (ctl *LotteryController) PostV1LotteryActivity(ctx jet.Ctx, req *req.LotteryActivityReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.LotteryService.AddOrUpdateActivity(ctx, req))
}

func (ctl *LotteryController) PostV1LotteryActivityList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	listPrize, count, err := ctl.LotteryService.ListActivity(ctx, params)
	pageResult := api.WrapPageResult(params, listPrize, count)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl *LotteryController) PostV1LotteryActivityStatus(ctx jet.Ctx, req *req.LotteryActivityStatusReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.LotteryService.UpdateActivityStatus(ctx, req))
}

func (ctl *LotteryController) PostV1LotteryActivityDel(ctx jet.Ctx, req *req.LotteryActivityReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "OK", ctl.LotteryService.DelActivity(ctx, req))
}

// =============================  lottery records  =======================================

func (ctl *LotteryController) PostV1LotteryRecordsList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	listPrize, count, err := ctl.LotteryService.ListLotteryRecords(ctx, params)
	pageResult := api.WrapPageResult(params, listPrize, count)
	return xjet.WrapperResult(ctx, pageResult, err)
}
