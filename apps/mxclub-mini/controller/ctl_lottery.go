package controller

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewLotteryController)
}

type LotteryController struct {
	jet.BaseJetController
	lotteryService *service.LotteryService
}

func NewLotteryController(LotteryService *service.LotteryService) jet.ControllerResult {
	return jet.NewJetController(&LotteryController{
		lotteryService: LotteryService,
	})
}

func (ctr *LotteryController) PostV1LotteryActivityList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	activity, total, err := ctr.lotteryService.ListLotteryPrize(ctx, params)
	if err != nil {
		return nil, err
	}
	pageResult := api.WrapPageResult(params, activity, total)
	return xjet.WrapperResult(ctx, pageResult, nil)
}

func (ctr *LotteryController) GetV1LotteryActivity0(ctx jet.Ctx, pathParam *api.PathParam) (*api.Response, error) {
	got, err := pathParam.GetInt64(0)
	if err != nil {
		return nil, errors.New("参数错误")
	}
	activityPrizeVO, err := ctr.lotteryService.FindActivityPrizeByActivityId(ctx, int(got))
	return xjet.WrapperResult(ctx, activityPrizeVO, err)
}
