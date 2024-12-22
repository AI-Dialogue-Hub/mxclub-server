package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
)

func (svc OrderService) Transfer(ctx jet.Ctx, req *req.TransferReq) error {
	executorTo := req.ExecutorTo
	if !req.IsValid {
		executorTo = -1
	}
	// 检查队友转单人是否在线
	if executorTo >= 0 {
		if !svc.userService.ExistsExecutor(ctx, executorTo) {
			ctx.Logger().Errorf("FindUserByDashId, not found executorTo:%v", executorTo)
			return errors.New("打手编号错误，指定打手不存在")
		}
		isOnline := svc.userService.CheckAssistantStatus(ctx, executorTo)
		inRunningOrder := svc.userService.CheckDasherInRunningOrder(ctx, executorTo)
		if !isOnline && inRunningOrder {
			ctx.Logger().Errorf("check failed, b1:%v, b2:%v, executorTo:%v", isOnline, inRunningOrder, executorTo)
			return errors.New("指定打手不在线或者正在进行中订单")
		}
	}
	userId := middleware.MustGetUserId(ctx)
	dasherPO, _ := svc.userService.FindUserById(ctx, userId)
	err := svc.transferRepo.InsertOne(&po.OrderTransfer{
		OrderId:      req.OrderId,
		ExecutorFrom: dasherPO.MemberNumber,
		ExecutorTo:   executorTo,
		Status:       enum.Transfer_PENDING,
	})
	if err != nil {
		ctx.Logger().Errorf("Transfer ERROR:%v", err)
		return errors.New("转单失败，请联系客服")
	}
	return nil
}

func (svc OrderService) RemoveTransferRecord(ctx jet.Ctx) error {
	userId := middleware.MustGetUserId(ctx)
	userPO, _ := svc.userService.FindUserById(ctx, userId)
	return svc.transferRepo.RemoveByDasherId(ctx, userPO.MemberNumber)
}
