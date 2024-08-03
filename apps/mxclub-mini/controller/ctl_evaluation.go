package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func (ctl OrderController) PutV1Evaluation(ctx jet.Ctx, req *req.EvaluationReq) (*api.Response, error) {
	err := ctl.orderService.AddEvaluation(ctx, req)
	return xjet.WrapperResult(ctx, "ok", err)
}
