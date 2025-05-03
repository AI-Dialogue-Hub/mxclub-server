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
	FindByOrderIds(ctx jet.Ctx, orderIds []string) (map[string][]*po.RewardRecord, error)
	// FindByOutTradeNo 打赏订单给到微信的唯一Id
	FindByOutTradeNo(ctx jet.Ctx, outTradeNo string) (*po.RewardRecord, error)
	ExistByOrderIdAndDasherNumber(ctx jet.Ctx, orderId string, dasherDBId uint) bool
	UpdateRewardStatus(ctx jet.Ctx, outTradeNo string, status enum.OrderStatus) error
	// AllRewardAmountByDasherId 查询所有打赏的钱，使用db id进行定位
	AllRewardAmountByDasherId(ctx jet.Ctx, dasherNumber uint) (float64, error)
	ListNoCountDuration(ctx jet.Ctx, startDateStr, endDateStr string, status enum.OrderStatus) ([]*po.RewardRecord, error)
	// ClearAllRewardByDasherId 清理打手所有打赏信息
	ClearAllRewardByDasherId(ctx jet.Ctx, dasherId uint) error
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
	// 通过用户db id查询所有打赏记录
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

func (repo RewardRepoImpl) FindByOrderIds(ctx jet.Ctx, orderIds []string) (map[string][]*po.RewardRecord, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(1, 100000)
	query.SetFilter("order_id in (?)", orderIds)
	rewardRecords, err := repo.ListNoCountByQuery(query)
	if err != nil {
		return nil, err
	}
	var m = make(map[string][]*po.RewardRecord)
	for _, record := range rewardRecords {
		if _, ok := m[record.OrderID]; ok {
			records := m[record.OrderID]
			records = append(records, record)
			m[record.OrderID] = records
		} else {
			records := make([]*po.RewardRecord, 0)
			records = append(records, record)
			m[record.OrderID] = records
		}
	}
	return m, nil
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

func (repo RewardRepoImpl) FindByOutTradeNo(ctx jet.Ctx, outTradeNo string) (*po.RewardRecord, error) {
	return repo.FindOne("out_trade_no = ?", outTradeNo)
}

func (repo RewardRepoImpl) ListNoCountDuration(ctx jet.Ctx, startDateStr, endDateStr string, status enum.OrderStatus) ([]*po.RewardRecord, error) {
	var (
		wrapper = xmysql.NewMysqlQuery()
		logger  = ctx.Logger()
	)
	wrapper.SetFilter("created_at >= ? and created_at <= ? and status = ?", startDateStr, endDateStr, status)
	wrapper.SetLimit(10000)
	rewardRecords, err := repo.ListNoCountByQuery(wrapper)
	if err != nil || rewardRecords == nil || len(rewardRecords) == 0 {
		logger.Errorf("cannot find any order, duration is: %v %v", startDateStr, endDateStr)
		return nil, err
	}
	return rewardRecords, nil
}

func (repo RewardRepoImpl) ClearAllRewardByDasherId(ctx jet.Ctx, dasherId uint) error {
	return repo.Remove("dasher_number = ?", dasherId)
}
