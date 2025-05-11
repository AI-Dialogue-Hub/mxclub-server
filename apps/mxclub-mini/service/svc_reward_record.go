package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/event"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewRewardRecordService)
	jet.Invoke(func(svc *RewardRecordService) {
		event.RegisterEvent("RewardService", event.EventRemoveDasher, svc.RemoveRewardRecord)
	})
}

func NewRewardRecordService(
	rewardRecordRepo repo.IRewardRecordRepo,
	userRepo userRepo.IUserRepo,
	wxpayCallbackRepo repo.IWxPayCallbackRepo,
	messageService *MessageService) *RewardRecordService {
	return &RewardRecordService{
		rewardRecordRepo:  rewardRecordRepo,
		userRepo:          userRepo,
		wxpayCallbackRepo: wxpayCallbackRepo,
		messageService:    messageService,
	}
}

type RewardRecordService struct {
	rewardRecordRepo  repo.IRewardRecordRepo
	userRepo          userRepo.IUserRepo
	wxpayCallbackRepo repo.IWxPayCallbackRepo
	messageService    *MessageService
}

// ============================================================================

func (svc RewardRecordService) AddByRewardStatus(
	ctx jet.Ctx,
	req *req.RewardPrepayReq,
	status enum.OrderStatus) (*po.RewardRecord, error) {

	var (
		userId = middleware.MustGetUserId(ctx)
		logger = ctx.Logger()
	)

	dasherPO, err := svc.userRepo.FindByMemberNumber(ctx, req.DasherId)

	if err != nil || dasherPO == nil || dasherPO.ID <= 0 {
		logger.Errorf("[AddByRewardStatus] dasherId:%v is invalid", req.DasherId)
		return nil, fmt.Errorf("指定打手:%v 不存在, 该打手可能已注销, 请联系客服", req.DasherId)
	}

	// 1. 检查是否已经有打赏记录，不允许重复打赏
	if svc.rewardRecordRepo.ExistByOrderIdAndDasherNumber(ctx, req.OrderId, dasherPO.ID) {
		return nil, fmt.Errorf("您已经有指定打手:%v的打赏了，请勿重复打赏哦", req.DasherId)
	}
	// 2. 添加打赏记录
	rewardPO := &po.RewardRecord{
		PurchaserID:  userId,
		OrderID:      req.OrderId,
		DasherID:     dasherPO.MemberNumber,
		DasherNumber: dasherPO.ID,
		DasherName:   dasherPO.Name,
		Remarks:      req.RewardNote,
		RewardAmount: req.RewardAmount,
		Status:       status,
		OutTradeNo:   wxpay.GenerateUniqueOrderNumber(),
	}

	err = svc.rewardRecordRepo.InsertOne(rewardPO)

	if err != nil {
		logger.Errorf("[AddByRewardStatus] add record failed, err:%v record info:%v",
			err, utils.ObjToJsonStr(rewardPO))
		return nil, fmt.Errorf("打赏失败, 请联系客服")
	}

	return rewardPO, nil
}

// PrePay 预支付，这里没有实现事务，可能有问题 todo@lfy
func (svc RewardRecordService) PrePay(ctx jet.Ctx, prepayReq *req.RewardPrepayReq) (*wxpay.PrePayDTO, error) {
	var (
		userId = middleware.MustGetUserId(ctx)
	)
	// 1. 插入一条预支付的打赏订单
	rewardRecord, err := svc.AddByRewardStatus(ctx, prepayReq, enum.PrePay)
	if err != nil {
		return nil, err
	}
	// 2. 生成微信 prepay 请求
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	prepayRequest := wxpay.NewPrepayRequest(prepayReq.RewardAmount, userPO.WxOpenId, rewardRecord.OutTradeNo)
	prePayDTO, err := wxpay.PrepayWithReward(ctx, prepayRequest)
	if err != nil {
		ctx.Logger().Errorf(
			"[WxPayService]prepay ERROR: %v\nprepayDTO:%v", err.Error(), utils.ObjToJsonStr(prePayDTO))
		return nil, errors.New("申请微信支付失败")
	}
	ctx.Logger().Infof(
		"用户: %v(%v) 打赏订单付款：%v，进行中，prepayDTO：%v",
		userPO, userPO.Name, prepayReq.RewardAmount, utils.ObjToJsonStr(prePayDTO))
	return prePayDTO, nil
}

func (svc RewardRecordService) WxpayNofity(ctx jet.Ctx, params *maps.LinkedHashMap[string, any]) {
	// 1. 解析回调参数
	transaction, err := wxpay.DecryptWxpayCallBack(ctx)
	if err != nil || transaction == nil {
		ctx.Logger().Errorf("[DecryptWxpayCallBack] %v", err)
		// 失败了，直接解析参数
		if transaction, err = wxpay.DecryptWxpayCallBackByParams(ctx, params); err != nil {
			ctx.Logger().Errorf("[DecryptWxpayCallBackByParams]ERROR %v", err)
			return
		}
	}

	var outTradeNo = *transaction.OutTradeNo

	callbackInfo, err := svc.wxpayCallbackRepo.FindByTraceNo(outTradeNo)

	// 幂等保护
	if err == nil && callbackInfo != nil && callbackInfo.ID > 0 {
		ctx.Logger().Errorf("duplicate callback, %v", utils.ObjToJsonStr(*transaction))
		return
	}

	// 2. 修改订单状态为支付成功
	_ = svc.PaySuccessOrder(ctx, outTradeNo)
	// 提示打手已经打赏了
	rewardPO, err := svc.rewardRecordRepo.FindByOutTradeNo(ctx, outTradeNo)

	if err == nil {
		// 2.1 给打手发消息
		_ = svc.messageService.PushSystemMessage(ctx, rewardPO.DasherNumber,
			fmt.Sprintf("打手您好，您的订单:%v 里老板为您进行打赏，打赏金额为:%v元",
				rewardPO.OrderID, utils.RoundToTwoDecimalPlaces(rewardPO.RewardAmount)),
		)
		// 2.2 给老板发消息
		_ = svc.messageService.PushSystemMessage(ctx, rewardPO.PurchaserID,
			fmt.Sprintf("老板您好，您的订单:%v，打赏打手%v(%v)已成功，打赏金额为:%v元",
				rewardPO.OrderID, rewardPO.DasherID,
				rewardPO.DasherName, utils.RoundToTwoDecimalPlaces(rewardPO.RewardAmount)),
		)
	} else {
		ctx.Logger().Errorf("FindByOutTradeNo error, outTradeNo:%v", outTradeNo)
	}

	// 3. 保存回调数据
	objToMap := utils.ObjToMap(*transaction)
	err = svc.wxpayCallbackRepo.InsertOne(&po.WxPayCallback{
		OutTradeNo: *transaction.OutTradeNo,
		RawData:    objToMap,
	})
	ctx.Logger().Infof("HandleWxpayNotify:%v", utils.ObjToJsonStr(objToMap))
	if err != nil {
		ctx.Logger().Errorf("[DecryptWxpayCallBack] %v", err)
	}
	return
}

func (svc RewardRecordService) PaySuccessOrder(ctx jet.Ctx, outTradeNo string) error {
	err := svc.rewardRecordRepo.UpdateRewardStatus(ctx, outTradeNo, enum.SUCCESS)
	if err != nil {
		ctx.Logger().Errorf("UpdateRewardStatus ERROR, %v", err)
		return err
	}
	return nil
}

func (svc RewardRecordService) List(ctx jet.Ctx, listReq *req.RewardListReq) ([]*vo.RewardVO, error) {
	queryWrapper := xmysql.NewMysqlQuery()
	queryWrapper.SetPage(listReq.Page, listReq.PageSize)
	if listReq.OrderId > 0 {
		queryWrapper.SetFilter("order_id = ?", listReq.OrderId)
	}
	rewardRecords, err := svc.rewardRecordRepo.ListNoCountByQuery(queryWrapper)
	if err != nil {
		ctx.Logger().Errorf("RewardRecordService#List ERROR, %v", err)
		return nil, errors.New("查询失败")
	}
	rewardVOS := utils.CopySlice[*po.RewardRecord, *vo.RewardVO](rewardRecords)
	return rewardVOS, nil
}

// RemoveRewardRecord 清理打手打赏信息
func (svc RewardRecordService) RemoveRewardRecord(ctx jet.Ctx) error {
	utils.RecoverAndLogError(ctx)
	userId := middleware.MustGetUserId(ctx)
	if err := svc.rewardRecordRepo.ClearAllRewardByDasherId(ctx, userId); err != nil {
		ctx.Logger().Errorf("RemoveRewardRecord ERROR, %v", err)
		return err
	}
	return nil
}
