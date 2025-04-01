package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewDashBoardController)
}

type DashBoardController struct {
	jet.BaseJetController
	svc *service.DashBoardService
}

func NewDashBoardController(svc *service.DashBoardService) jet.ControllerResult {
	return jet.NewJetController(&DashBoardController{svc: svc})
}

func (ctl DashBoardController) GetV1Dashboard(ctx jet.Ctx, req *req.SaleForDayReq) (*api.Response, error) {
	saleByDuration, err := ctl.svc.FindSaleByDuration(req)
	return xjet.WrapperResult(ctx, saleByDuration, err)
}
