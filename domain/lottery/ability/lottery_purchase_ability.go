package ability

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
	"time"
)

func (ability *LotteryAbility) FindPurchaseRecord(
	ctx jet.Ctx, userId uint, activityId uint, canDrawLottery bool,
) ([]*po.LotteryPurchaseRecord, error) {
	query := xmysql.NewMysqlQuery()
	query.SetFilter("user_id = ?", userId)
	query.SetFilter("activity_id = ?", activityId)
	query.SetFilter("lottery_qualified = ?", true)
	if canDrawLottery {
		query.SetFilter("lottery_used = ?", false)
		purchaseRecords, err := ability.lotteryPurchaseRecordsRepo.FindOneByWrapper(query)
		if err != nil || purchaseRecords == nil || purchaseRecords.ID <= 0 {
			ctx.Logger().Errorf("ERROR:%v, purchaseRecords:%+v", err, purchaseRecords)
			return nil, errors.Wrap(err, "find purchase records error")
		}
		return utils.ToSlice(purchaseRecords), nil
	} else {
		// 已经获得抽奖资格未使用 和 未获得抽奖资格，状态是hold的订单
		query.SetFilter(
			"((lottery_used = ? and lottery_qualified = ?) or (lottery_used = ? and lottery_qualified = ?))",
			true, true, false, false)
		purchaseRecords, err := ability.lotteryPurchaseRecordsRepo.FindByWrapper(query)
		if err != nil || purchaseRecords == nil || len(purchaseRecords) <= 0 {
			ctx.Logger().Errorf("ERROR:%v", err)
			return nil, errors.Wrap(err, "find purchase records error")
		}
		return utils.Filter(purchaseRecords, func(in *po.LotteryPurchaseRecord) bool { return in.ID > 0 }), nil
	}
}
func (ability *LotteryAbility) FindPurchaseRecordByTransactionId(
	ctx jet.Ctx, transactionId string) (*po.LotteryPurchaseRecord, error) {
	one, err := ability.lotteryPurchaseRecordsRepo.FindOne("transaction_id = ?", transactionId)
	if err != nil {
		ctx.Logger().Errorf("find purchase record err:%v", err)
		return nil, errors.Wrap(err, "find purchase record error")
	}
	return one, nil
}

func (ability *LotteryAbility) AddPurchaseRecordByActivityId(
	ctx jet.Ctx, userId uint, activityId uint, transactionNo string) (*po.LotteryPurchaseRecord, error) {
	// 0. 检查活动是否存在
	a, err := ability.lotteryActivityRepo.FindByID(activityId)
	if err != nil || a == nil || a.ID <= 0 {
		ctx.Logger().Errorf(
			"[LotteryAbility#AddPurchaseRecordByActivityId] cannot find activityId:%v", activityId)
		return nil, errors.Wrap(err, "Activity not found")
	}
	purchaseRecordPO := &po.LotteryPurchaseRecord{
		UserID:           userId,
		ActivityID:       activityId,
		TransactionID:    transactionNo,
		PurchaseAmount:   a.ActivityPrice,
		PurchaseTime:     time.Now(),
		PurchaseStatus:   enum.PurchaseStatusHold,
		PaymentMethod:    enum.PaymentMethodWeChat,
		LotteryQualified: false,
		LotteryUsed:      false,
		IPAddress:        utils.Ptr(xjet.GetClientIP(ctx)),
		DeviceInfo:       nil,
	}
	err = ability.lotteryPurchaseRecordsRepo.InsertOne(purchaseRecordPO)
	if err != nil {
		ctx.Logger().Errorf("[LotteryAbility#AddPurchaseRecordByActivityId] err:%v", err)
		return nil, errors.Wrap(err, "insert purchase record error")
	}
	ctx.Logger().Infof("[LotteryAbility#AddPurchaseRecordByActivityId] purchase record:%v", purchaseRecordPO)
	return purchaseRecordPO, nil
}

func (ability *LotteryAbility) UpdatePurchaseRecordStatus(
	ctx jet.Ctx, transactionId string, status enum.PurchaseStatusEnum) error {
	if status != enum.PurchaseStatusSuccess {
		ctx.Logger().Errorf("LotteryAbility status is not success, transactionId=%v", transactionId)
		return errors.New("cannot implement the status")
	}
	updateMap := map[string]any{
		"lottery_qualified": true,
		"lottery_used":      false,
		"payment_time":      time.Now(),
	}
	return ability.lotteryPurchaseRecordsRepo.Update(updateMap, "transaction_id = ?", transactionId)
}
