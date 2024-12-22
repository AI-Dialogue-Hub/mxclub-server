package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/middleware"
)

func (svc OrderService) RemoveWithdrawalRecord(ctx jet.Ctx) error {
	return svc.withdrawalRepo.RemoveWithdrawalRecord(ctx, middleware.MustGetUserId(ctx))
}
