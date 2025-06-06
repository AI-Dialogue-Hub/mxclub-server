package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/utils"
)

func (svc *OrderService) ListDeduction(ctx jet.Ctx, listReq *req.DeductionListReq) ([]*vo.DeductionVO, int64, error) {
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

	vos := utils.Map[*po.Deduction, *vo.DeductionVO](listDeduction, func(in *po.Deduction) *vo.DeductionVO {
		var userInfo string
		if userPO, err := svc.userRepo.FindByIdAroundCache(ctx, in.UserID); err == nil {
			userInfo = fmt.Sprintf("%v(%v)", userPO.MemberNumber, userPO.Name)
		}
		return &vo.DeductionVO{
			ID:              in.ID,
			UserID:          in.UserID,
			UserInfo:        userInfo,
			ConfirmPersonId: in.ConfirmPersonId,
			Amount:          in.Amount,
			Reason:          in.Reason,
			Status:          in.Status.DisPlayName(),
			CreatedAt:       in.CreatedAt,
		}
	})
	return vos, total, nil
}

func (svc *OrderService) AddDeduction(ctx jet.Ctx, addReq *req.DeductionAddReq) error {
	deductPO := utils.MustCopy[po.Deduction](addReq)
	err := svc.deductionRepo.InsertOne(deductPO)
	if err != nil {
		ctx.Logger().Errorf("[orderService]AddDeduction ERROR:%v", err.Error())
		return errors.New("添加失败")
	}
	return nil
}

func (svc *OrderService) UpdateDeduction(ctx jet.Ctx, updateReq *req.DeductionUpdateReq) error {
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
	// 如果是评星，要删除评星记录
	if updateReq.Status == "已拒绝" || updateReq.Status == enum.Deduct_REJECT.DisPlayName() {
		if deductionPO, _ := svc.deductionRepo.FindByID(updateReq.ID); deductionPO != nil && deductionPO.ID > 0 {
			orderNo := deductionPO.OrderNo
			dasherId := deductionPO.DasherId
			// TODO@lfy 这里其实没有区分罚款的类型
			_ = svc.evaluationRepo.RemoveByOrderIdAndDasherId(ctx, orderNo, dasherId)
		}
	}
	return nil
}

func (svc *OrderService) Delete(ctx jet.Ctx, id uint) error {
	err := svc.deductionRepo.RemoveByID(id)
	if err != nil {
		ctx.Logger().Errorf("[orderService]Delete ERROR:%v", err.Error())
		return errors.New("删除失败")
	}
	return nil
}
