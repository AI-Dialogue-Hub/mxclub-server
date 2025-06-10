package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	constantMini "mxclub/apps/mxclub-mini/entity/constant"
	"mxclub/apps/mxclub-mini/middleware"
)

func (svc *OrderService) RemoveWithdrawalRecord(ctx jet.Ctx) error {
	if value, exists := ctx.Get(constantMini.LOGOUT_DASHER_ID); exists {
		dasherId := value.(int)
		return svc.withdrawalRepo.RemoveWithdrawalRecordByDasherId(ctx, dasherId)
	} else {
		userId := middleware.MustGetUserId(ctx)
		ctx.Logger().Infof("[OrderService#RemoveWithdrawalRecord] userId:%v", userId)
		return svc.withdrawalRepo.RemoveWithdrawalRecord(ctx, userId)
	}
}
