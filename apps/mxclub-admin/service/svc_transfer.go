package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/common/xmysql"
)

// ClearAllDasherInfo 清空所有打手信息，重新派单到大厅
func (svc OrderService) ClearAllDasherInfo(ctx jet.Ctx, id uint) error {
	err := svc.orderRepo.ClearOrderDasherInfo(ctx, id)
	if err != nil {
		ctx.Logger().Errorf("[ClearAllDasherInfo]err:%v", err)
		return errors.New("转单失败")
	}
	return nil
}

func (svc OrderService) ListTransferInfo(ctx jet.Ctx, params *req.TransferListReq) ([]*vo.TransferVO, int64, error) {
	query := new(xmysql.MysqlQuery)
	query.SetPage(params.Page, params.PageSize)
	if params.Status >= 0 {
		query.SetFilter("status = ?", params.Status)
	}
	transfers, count, err := svc.transferRepo.ListByWrapper(ctx, query)
	if err != nil {
		ctx.Logger().Errorf("ListByWrapper ERROR:%v", err)
		return nil, 0, errors.New("查找失败")
	}
	var vos []*vo.TransferVO
	for _, transferPO := range transfers {
		p1, e1 := svc.userRepo.FindByMemberNumber(ctx, transferPO.ExecutorFrom)
		if e1 != nil {
			ctx.Logger().Errorf("FindByMemberNumber ERROR:%v  p1:%v, transferPO:%v", e1, p1, transferPO)
			continue
		}
		var p2Name string
		if transferPO.ExecutorTo >= 0 {
			p2, e2 := svc.userRepo.FindByMemberNumber(ctx, transferPO.ExecutorTo)
			if e2 != nil && p2 != nil {
				p2Name = p2.Name
			}
		}
		vos = append(vos, &vo.TransferVO{
			ID:               transferPO.ID,
			OrderId:          transferPO.OrderId,
			ExecutorFrom:     transferPO.ExecutorFrom,
			ExecutorFromName: p1.Name,
			ExecutorTo:       transferPO.ExecutorTo,
			ExecutorToName:   p2Name,
			Status:           transferPO.Status,
			CreatedAt:        transferPO.CreatedAt,
		})
	}
	return vos, count, nil
}

func (svc OrderService) RemoveTransfer(ctx jet.Ctx, id int64) error {
	err := svc.transferRepo.RemoveByID(id)
	if err != nil {
		ctx.Logger().Errorf("RemoveByID ERROR:%v", err)
		return errors.New("删除失败")
	}
	return nil
}

func (svc OrderService) UpdateTransfer(ctx jet.Ctx, transferVO *vo.TransferVO) error {
	oldInfo, _ := svc.transferRepo.FindByID(transferVO.ID)
	if oldInfo.Status != enum.Transfer_PENDING {
		return errors.New("只能修改待处理的申请")
	}
	fromDasherPO, _ := svc.userRepo.FindByMemberNumber(ctx, transferVO.ExecutorFrom)

	if oldInfo.Status != transferVO.Status {
		if transferVO.Status == enum.Transfer_SUCCES {
			// 通过申请，查找是否指定打手
			executorTo := transferVO.ExecutorTo
			if executorTo >= 0 {
				if online := svc.userRepo.CheckAssistantStatus(ctx, executorTo); !online {
					return errors.New(fmt.Sprintf("指定打手不在线，打手Id为:%v", executorTo))
				}
				if inRunningOrder := svc.CheckDasherInRunningOrder(ctx, executorTo); inRunningOrder {
					return errors.New(fmt.Sprintf("指定打手正在进行中订单，打手Id为:%v", executorTo))
				}
				svc.transferRepo.UpdateById(map[string]any{"status": enum.Transfer_SUCCES}, transferVO.ID)
				// 1. 更新订单
				toDasherPO, _ := svc.userRepo.FindByMemberNumber(ctx, transferVO.ExecutorTo)
				orderUpdate := map[string]any{"executor_id": toDasherPO.MemberNumber, "executor_name": toDasherPO.Name}
				if err := svc.orderRepo.UpdateById(orderUpdate, transferVO.OrderId); err != nil {
					ctx.Logger().Errorf("UpdateById ERROR:%v", err)
					return errors.New("订单更新异常，转单失败")
				}
				// 2. 发送消息提示
				svc.messageService.PushSystemMessage(
					ctx,
					toDasherPO.ID,
					fmt.Sprintf("您被打手:%v(%v)指定转单，请前往订单列表进行查看", fromDasherPO.Name, fromDasherPO.MemberNumber),
				)
				svc.messageService.PushSystemMessage(ctx, fromDasherPO.ID, "转单申请已通过")
			} else {
				// 清理所有打手信息，重新发往大厅
				if err := svc.ClearAllDasherInfo(ctx, uint(transferVO.OrderId)); err != nil {
					return err
				}
				// 更新转单订单
				svc.transferRepo.UpdateById(map[string]any{"status": enum.Transfer_SUCCES}, transferVO.ID)
			}
		} else if transferVO.Status == enum.Transfer_REJECT {
			if err := svc.transferRepo.UpdateById(map[string]any{"status": enum.Transfer_REJECT}, transferVO.ID); err != nil {
				ctx.Logger().Errorf("UpdateById ERROR:%v", err)
				return errors.New("转单更新失败")
			}
		}
	}
	// 指定打手修改
	if oldInfo.ExecutorTo != transferVO.ExecutorTo {
		executorTo := transferVO.ExecutorTo
		if online := svc.userRepo.CheckAssistantStatus(ctx, executorTo); !online {
			return errors.New(fmt.Sprintf("指定打手不在线，打手Id为:%v", executorTo))
		}
		if inRunningOrder := svc.CheckDasherInRunningOrder(ctx, executorTo); inRunningOrder {
			return errors.New(fmt.Sprintf("指定打手正在进行中订单，打手Id为:%v", executorTo))
		}
		if err := svc.transferRepo.UpdateById(map[string]any{"executor_to": executorTo}, transferVO.ID); err != nil {
			ctx.Logger().Errorf("UpdateById ERROR:%v", err)
			return errors.New("转单更新失败")
		}
	}
	return nil
}
