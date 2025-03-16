package controller

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewRewardController)
}

type RewardController struct {
	jet.BaseJetController
	svc *service.RewardRecordService
}

func NewRewardController(svc *service.RewardRecordService) jet.ControllerResult {
	return jet.NewJetController(&RewardController{
		svc: svc,
	})
}

// =======================================================================

func (ctl RewardController) PostV1RewardList(ctx jet.Ctx, req *req.RewardListReq) (*api.Response, error) {
	ctx.Logger().Infof("[PostV1RewardPrepay] req => %v", utils.ObjToJsonStr(req))
	rewardVOList, err := ctl.svc.List(ctx, req)
	return xjet.WrapperResult(ctx, api.WrapPageResult(req.PageParams, rewardVOList, 0), err)
}

func (ctl RewardController) PostV1RewardPrepay(ctx jet.Ctx, req *req.RewardPrepayReq) (*api.Response, error) {
	ctx.Logger().Infof("[PostV1RewardPrepay] req => %v", utils.ObjToJsonStr(req))
	prePayDTO, err := ctl.svc.PrePay(ctx, req)
	return xjet.WrapperResult(ctx, prePayDTO, err)
}

// ======= 回调 =========

func (ctl RewardController) GetV1RewardWxpayNotify(
	ctx jet.Ctx,
	params *maps.LinkedHashMap[string, any]) (*api.Response, error) {

	go ctl.svc.WxpayNofity(ctx, params)

	return xjet.WrapperResult(ctx, "ok", nil)
}

func (ctl RewardController) PostV1RewardWxpayNotify(
	ctx jet.Ctx,
	params *maps.LinkedHashMap[string, any]) (*api.Response, error) {

	go ctl.svc.WxpayNofity(ctx, params)

	return xjet.WrapperResult(ctx, "ok", nil)
}
