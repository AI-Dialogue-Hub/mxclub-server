package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/repo"
)

func init() {
	jet.Provide(NewDashBoardService)
}

type DashBoardService struct {
	orderRepo repo.IOrderRepo
}

func NewDashBoardService(orderRepo repo.IOrderRepo) *DashBoardService {
	return &DashBoardService{orderRepo: orderRepo}
}

func (svc *DashBoardService) FindSaleByDuration(req *req.SaleForDayReq) (*vo.DashBoardVO, error) {

	return nil, nil
}
