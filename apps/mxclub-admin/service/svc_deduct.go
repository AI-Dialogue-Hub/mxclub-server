package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/utils"
)

func (svc OrderService) ListDeduction(ctx jet.Ctx, listReq *req.DeductionListReq) ([]*vo.DeductionVO, int64, error) {
	d := &dto.DeductionDTO{
		PageParams: listReq.PageParams,
		Ge:         listReq.Ge,
		Le:         listReq.Le,
		Status:     utils.CaseToPoint(enum.Of(listReq.Status)),
	}
	listDeduction, total, err := svc.deductionRepo.ListDeduction(ctx, d)
	if err != nil {
		ctx.Logger().Errorf("[orderService]ListWithdraw ERROR:%v", err.Error())
		return nil, 0, errors.New("获取失败")
	}
	vos := utils.CopySlice[*po.Deduction, *vo.DeductionVO](listDeduction)
	utils.ForEach(vos, func(ele *vo.DeductionVO) {
		ele.Status = enum.DeductStatus(ele.Status).DisPlayName()
	})
	return vos, total, nil
}

func (svc OrderService) AddDeduction(ctx jet.Ctx, addReq *req.DeductionAddReq) error {
	deductPO := utils.MustCopy[po.Deduction](addReq)
	err := svc.deductionRepo.InsertOne(deductPO)
	if err != nil {
		ctx.Logger().Errorf("[orderService]AddDeduction ERROR:%v", err.Error())
		return errors.New("添加失败")
	}
	return nil
}

func (svc OrderService) UpdateDeduction(ctx jet.Ctx, updateReq *req.DeductionUpdateReq) error {
	updateMap := map[string]any{
		"amount": updateReq.Amount,
		"reason": updateReq.Reason,
		"status": enum.Of(updateReq.Status),
	}
	err := svc.deductionRepo.Update(updateMap, "id = ?", updateReq.ID)
	if err != nil {
		ctx.Logger().Errorf("[orderService]UpdateDeduction ERROR:%v", err.Error())
		return errors.New("修改失败")
	}
	return nil
}

func (svc OrderService) Delete(ctx jet.Ctx, id uint) error {
	err := svc.deductionRepo.RemoveByID(id)
	if err != nil {
		ctx.Logger().Errorf("[orderService]Delete ERROR:%v", err.Error())
		return errors.New("删除失败")
	}
	return nil
}
