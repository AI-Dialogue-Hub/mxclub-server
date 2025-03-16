package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewRewardRecordRepo)
}

type IRewardRecordRepo interface {
	xmysql.IBaseRepo[po.RewardRecord]
	FindByOrderIdAndDasherNumber(ctx jet.Ctx, orderId string, dasherDBId uint) (*po.RewardRecord, error)
	ExistByOrderIdAndDasherNumber(ctx jet.Ctx, orderId string, dasherDBId uint) bool
	UpdateRewardStatus(ctx jet.Ctx, outTradeNo string, status enum.OrderStatus) error
	// AllRewardAmountByDasherId 查询所有打赏的钱，使用db id进行定位
	AllRewardAmountByDasherId(ctx jet.Ctx, dasherNumber uint) (float64, error)
}

func NewRewardRecordRepo(db *gorm.DB) IRewardRecordRepo {
	repo := new(RewardRepoImpl)
	repo.SetDB(db)
	repo.ModelPO = new(po.RewardRecord)
	repo.Ctx = context.Background()
	return repo
}

type RewardRepoImpl struct {
	xmysql.BaseRepo[po.RewardRecord]
}

// ====================================================================

const (
	sqlAllRewardAmountByDasherId = `SELECT IFNULL(SUM(reward_amount), 0) AS rewardAmount 
                                       FROM reward_records 
                                       WHERE dasher_number = ? AND deleted_at IS NULL`
)

func (repo RewardRepoImpl) FindByOrderIdAndDasherNumber(
	ctx jet.Ctx,
	orderId string,
	dasherDBId uint) (*po.RewardRecord, error) {
	return repo.FindOne("order_id = ? and dasher_number = ?", orderId, dasherDBId)
}

func (repo RewardRepoImpl) ExistByOrderIdAndDasherNumber(
	ctx jet.Ctx,
	orderId string,
	dasherDBId uint) bool {

	rewardRecordPO, err := repo.FindByOrderIdAndDasherNumber(ctx, orderId, dasherDBId)

	if err == nil && rewardRecordPO != nil && rewardRecordPO.ID > 0 {
		return true
	}

	return false
}

func (repo RewardRepoImpl) UpdateRewardStatus(ctx jet.Ctx, outTradeNo string, status enum.OrderStatus) error {
	updateWrapper := xmysql.NewMysqlUpdate()
	updateWrapper.SetFilter("out_trade_no = ?", outTradeNo)
	updateWrapper.Set("status", status)
	return repo.UpdateByWrapper(updateWrapper)
}

func (repo RewardRepoImpl) AllRewardAmountByDasherId(ctx jet.Ctx, dasherNumber uint) (float64, error) {
	var rewardAmount float64
	err := repo.DB().Raw(sqlAllRewardAmountByDasherId, dasherNumber).Scan(&rewardAmount).Error
	if err != nil {
		ctx.Logger().Errorf("AllRewardAmountByDasherId ERROR, %v", err)
		return 0, err
	}
	return rewardAmount, nil
}
