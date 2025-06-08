package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/middleware"
)

func (svc *OrderService) RemoveWithdrawalRecord(ctx jet.Ctx) error {
	userId := middleware.MustGetUserId(ctx)
	ctx.Logger().Infof("[OrderService#RemoveWithdrawalRecord] userId:%v", userId)
	return svc.withdrawalRepo.RemoveWithdrawalRecord(ctx, userId)
}
