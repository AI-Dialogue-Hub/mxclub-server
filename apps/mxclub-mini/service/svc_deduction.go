package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/constant"
	"mxclub/pkg/utils"
	"strings"
	"time"
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
		ctx.Logger().Errorf("[orderService]ListWithdraw ERROR:%v", err.Error())
		return nil, errors.New("获取失败")
	}
	vos := utils.CopySlice[*po.Deduction, *vo.DeductionVO](listDeduction)
	utils.ForEach(vos, func(ele *vo.DeductionVO) {
		ele.Status = enum.DeductStatus(ele.Status).DisPlayName()
	})
	return vos, nil
}

var (
	syncDeductionInfoLogger = xlog.NewWith("syncDeductionInfo")
	deductionDDL            = constant.Duration_2_Day
)

// SyncDeductionInfo 同步处罚情况
// 1. 如果处罚记录超过五天，将处罚状态标记为正式处罚
// 2. 如果处罚没满五天，在第三天的时候，再次提醒用户一次，找客服解除处罚
func (svc OrderService) SyncDeductionInfo() {
	defer utils.RecoverByPrefix(syncDeductionInfoLogger, "syncDeductionInfo")
	// 1. 处罚记录超过五天的
	deductions, err := svc.deductionRepo.FindDeDuctListBeyondDuration(deductionDDL)
	if err != nil || len(deductions) <= 0 {
		syncDeductionInfoLogger.Errorf("FindDeDuctListBeyondDuration ERROR:%v, deductions: %v", err, deductions)
		return
	}
	// userId -> PO
	userId2POMap := utils.SliceToMap(deductions, func(ele *po.Deduction) uint { return ele.UserID })
	userId2POMap.ForEach(func(userId uint, deductionList []*po.Deduction) {
		// 1.1 批量更新处罚记录
		ids := utils.Map[*po.Deduction, uint](deductions, func(in *po.Deduction) uint {
			return in.ID
		})
		err = svc.deductionRepo.UpdateStatusByIds(ids, enum.Deduct_SUCCESS)
		if err != nil {
			syncDeductionInfoLogger.Errorf("UpdateStatusByIds ERROR:%v", err)
		} else {
			// 1.2 给用户发送消息，提示被处罚
			builder := new(strings.Builder)
			builder.WriteString("尊敬的打手您好，您有以下处罚内容已超过2天未申述，系统已经进行处罚：\n")
			for _, deduction := range deductionList {
				builder.WriteString(fmt.Sprintf("处罚Id：%v，处罚原因：%v\n", deduction.ID, deduction.Reason))
			}
			_ = svc.messageService.PushSystemMessage(xjet.NewDefaultJetContext(), userId, builder.String())
		}
	})

	// 2. 处罚记录超过1天但是还没超过2天的
	deductions, err = svc.deductionRepo.FindDeDuctListWithDurations(constant.Duration_1_Day, constant.Duration_2_Day)
	if err != nil {
		syncDeductionInfoLogger.Errorf("FindDeDuctListWithDurations ERROR:%v", err)
		return
	}
	// userId -> PO
	userId2POMap = utils.SliceToMap(deductions, func(ele *po.Deduction) uint { return ele.UserID })
	userId2POMap.ForEach(func(userId uint, deductionList []*po.Deduction) {
		// 2.1 给用户发送消息，提示被处罚
		builder := new(strings.Builder)
		builder.WriteString("尊敬的打手您好，您有以下处罚内容已超过一天，如有异议请即使找客服进行申述，超过两天未申述，系统将进行处罚：\n")
		for _, deduction := range deductionList {
			builder.WriteString(fmt.Sprintf("处罚Id：%v，处罚原因：%v\n", deduction.ID, deduction.Reason))
		}
		syncDeductionInfoLogger.Infof("push DeductionInfo, userId:%v, message: %v", userId, builder.String())
		_ = svc.messageService.PushSystemMessage(xjet.NewDefaultJetContext(), userId, builder.String())
	})
}

func (svc OrderService) SyncPrePayOrder() {
	var (
		now      = time.Now()
		duration = -time.Second * 60 * 8 // 八分钟还没支付就移除掉
	)
	err := svc.orderRepo.Remove("created_at <= ? and order_status = ?", now.Add(duration), enum.PrePay)
	if err != nil {
		xlog.Errorf("SyncPrePayOrder ERROR:%v", err)
	}
}

func (svc OrderService) RemoveDeductRecord(ctx jet.Ctx) error {
	return svc.deductionRepo.RemoveDasher(ctx, middleware.MustGetUserId(ctx))
}
