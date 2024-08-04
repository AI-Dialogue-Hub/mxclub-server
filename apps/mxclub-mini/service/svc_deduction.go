package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/utils"
)

func (svc OrderService) ListDeduction(ctx jet.Ctx, listReq *req.DeductionListReq) ([]*vo.DeductionVO, error) {
	d := &dto.DeductionDTO{
		PageParams: listReq.PageParams,
		Ge:         listReq.Ge,
		Le:         listReq.Le,
		UserId:     middleware.MustGetUserId(ctx),
	}
	listDeduction, _, err := svc.deductionRepo.ListDeduction(ctx, d)
	if err != nil {
		ctx.Logger().Errorf("[OrderService]ListWithdraw ERROR:%v", err.Error())
		return nil, errors.New("获取失败")
	}
	vos := utils.CopySlice[*po.Deduction, *vo.DeductionVO](listDeduction)
	utils.ForEach(vos, func(ele *vo.DeductionVO) {
		ele.Status = enum.DeductStatus(ele.Status).DisPlayName()
	})
	return vos, nil
}
