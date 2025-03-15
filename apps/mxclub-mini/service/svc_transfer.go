package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/constant"
)

func (svc OrderService) Transfer(ctx jet.Ctx, req *req.TransferReq) error {
	executorTo := req.ExecutorTo
	if !req.IsValid {
		executorTo = -1
	}
	if err := checkTransferStatus(ctx, req.OrderId, executorTo); err != nil {
		return err
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
	_ = svc.messageService.PushSystemMessage(ctx, userId, "转单申请已提交，请联系董事长同意")
	return nil
}

// 检查转单申请是否重复提交
func checkTransferStatus(ctx jet.Ctx, orderId uint64, executorId int) error {
	// 防抖
	debounceKey := fmt.Sprintf("transfer_orderId_%v_executorId_%v", orderId, executorId)
	if err := xredis.Debounce(debounceKey, constant.Duration_minute); err != nil {
		ctx.Logger().Errorf("duplicated evaluation,evaluationReq => %v", debounceKey)
		return errors.New("请勿多次重复提交转单申请，请联系董事长进行审核")
	}
	return nil
}

func (svc OrderService) RemoveTransferRecord(ctx jet.Ctx) error {
	userId := middleware.MustGetUserId(ctx)
	userPO, _ := svc.userService.FindUserById(ctx, userId)
	return svc.transferRepo.RemoveByDasherId(ctx, userPO.MemberNumber)
}
