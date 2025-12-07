package ability

import (
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
	"slices"

	"github.com/fengyuan-liang/GoKit/collection/stream"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func init() {
	jet.Provide(NewLotteryActivity)
}

// ILotteryAbility 是领域根，不是仓储实现
type ILotteryAbility interface {
	FetchLotteryPrizeType() *dto.LotteryPrizeTypeDTO
	AddPrize(ctx jet.Ctx, activityId uint, lotteryPrize *po.LotteryPrize) error
	DelPrize(ctx jet.Ctx, lotteryPrizeId uint) error
	// ===========
	AddOrUpdateActivity(ctx jet.Ctx, lotteryAbility *po.LotteryActivity) error
	ListActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error)
	ListHotActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error)
	IncrementSalesVolume(ctx jet.Ctx, activityId uint, count int) error
	ListActivityPrize(ctx jet.Ctx, params *api.PageParams) ([]*dto.LotteryActivityDTO, int64, error)
	DelActivity(ctx jet.Ctx, lotteryActivityId uint) error
	FindActivityPrizeByActivityId(ctx jet.Ctx, activityId uint) (*dto.LotteryActivityDTO, error)
	AddLotteryRecords(ctx jet.Ctx, lotteryRecords *po.LotteryRecords) error
	// FindFallbackPrize 查找多次不中后，兜底的二等奖
	FindFallbackPrize(ctx jet.Ctx, activityId uint) (*po.LotteryPrize, error)
	// FindPurchaseRecord 查看用户的抽奖购买记录
	FindPurchaseRecord(
		ctx jet.Ctx, userId uint, activityId uint, canDrawLottery bool) ([]*po.LotteryPurchaseRecord, error)
	FindPurchaseRecordByTransactionId(
		ctx jet.Ctx, transactionId string) (*po.LotteryPurchaseRecord, error)
	AddPurchaseRecordByActivityId(
		ctx jet.Ctx, userId uint, activityId uint, transactionNo string) (*po.LotteryPurchaseRecord, error)
	PayAndUpdatePurchaseRecordStatus(
		ctx jet.Ctx, transactionId string, status enum.PurchaseStatusEnum) error
	// UpdatePurchaseRecordStatus 更新购买记录状态
	//
	// 	- @param lotteryQualifiedStatus 是否获得抽奖资格
	// 	- @param lotteryUsedStatus 是否已使用抽奖资格
	UpdatePurchaseRecordStatus(
		ctx jet.Ctx, transactionId string, lotteryQualifiedStatus, lotteryUsedStatus bool) error
}

type LotteryAbility struct {
	lotteryRepo                repo.ILotteryRepo
	lotteryPrizeRepo           repo.ILotteryPrizeRepo
	lotteryActivityRepo        repo.ILotteryActivityRepo
	relationRepo               repo.ILotteryActivityPrizeRelationRepo
	lotteryRecordsRepo         repo.ILotteryRecordsRepo
	lotteryPurchaseRecordsRepo repo.ILotteryPurchaseRecordsRepo
}

func NewLotteryActivity(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivityRepo repo.ILotteryActivityRepo,
	lotteryRepo repo.ILotteryRepo,
	relationRepo repo.ILotteryActivityPrizeRelationRepo,
	lotteryRecordsRepo repo.ILotteryRecordsRepo,
	lotteryPurchaseRecordsRepo repo.ILotteryPurchaseRecordsRepo) ILotteryAbility {
	lotteryAbilityInstance = &LotteryAbility{
		lotteryRepo:                lotteryRepo,
		lotteryPrizeRepo:           lotteryPrizeRepo,
		lotteryActivityRepo:        lotteryActivityRepo,
		relationRepo:               relationRepo,
		lotteryRecordsRepo:         lotteryRecordsRepo,
		lotteryPurchaseRecordsRepo: lotteryPurchaseRecordsRepo,
	}
	return lotteryAbilityInstance
}

var lotteryAbilityInstance ILotteryAbility

func FetchLotteryAbilityInstance() ILotteryAbility {
	return lotteryAbilityInstance
}

func (ability *LotteryAbility) FetchLotteryPrizeType() *dto.LotteryPrizeTypeDTO {
	return &dto.LotteryPrizeTypeDTO{PrizeType: enum.PrizeTypeNames}
}

func (ability *LotteryAbility) AddPrize(ctx jet.Ctx, activityId uint, lotteryPrize *po.LotteryPrize) error {
	// 1. 插入奖品
	if err := ability.lotteryPrizeRepo.InsertOne(lotteryPrize); err != nil {
		return err
	}
	ctx.Logger().Infof("[LotteryAbility#AddPrize] success, lotteryPrize=%v", utils.ObjToJsonStr(lotteryPrize))
	// 2. 插入活动-奖品关系
	if err := ability.relationRepo.InsertOne(&po.ActivityPrizeRelation{
		ActivityID: activityId,
		PrizeID:    lotteryPrize.ID,
	}); err != nil {
		return err
	}
	return nil
}

func (ability *LotteryAbility) DelPrize(ctx jet.Ctx, lotteryPrizeId uint) error {
	// 1. 删除奖品
	if err := ability.lotteryPrizeRepo.RemoveByID(lotteryPrizeId); err != nil {
		ctx.Logger().Errorf("[LotteryAbility#DelPrize] err:%v", err)
		return err
	}
	// 2. 删除跟活动的关联关系
	if err := ability.relationRepo.DelByPrizeId(ctx, lotteryPrizeId); err != nil {
		ctx.Logger().Errorf("[LotteryAbility#DelPrize] err:%v", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	return nil
}

func (ability *LotteryAbility) AddOrUpdateActivity(ctx jet.Ctx, lotteryAbility *po.LotteryActivity) error {
	if lotteryAbility.ID > 0 {
		if err := ability.lotteryActivityRepo.UpdateById(lotteryAbility, lotteryAbility.ID); err != nil {
			ctx.Logger().Errorf("[LotteryAbility#AddOrUpdateActivity] err:%v", err)
			return err
		}
		return nil
	}
	return ability.lotteryActivityRepo.InsertOne(lotteryAbility)
}

func (ability *LotteryAbility) ListActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error) {
	lotteryActivities, lotteryActivityCount, err := ability.lotteryActivityRepo.List(params.Page, params.PageSize, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "ListActivity error")
	}
	if lotteryActivities == nil || len(lotteryActivities) == 0 {
		return make([]*po.LotteryActivity, 0, 0), lotteryActivityCount, nil
	}
	return lotteryActivities, lotteryActivityCount, nil
}

func (ability *LotteryAbility) ListHotActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error) {
	lotteryActivities, lotteryActivityCount, err := ability.
		lotteryActivityRepo.List(params.Page, params.PageSize, "is_hot", true)
	if err != nil {
		return nil, 0, errors.Wrap(err, "ListActivity error")
	}
	if lotteryActivities == nil || len(lotteryActivities) == 0 {
		return make([]*po.LotteryActivity, 0, 0), lotteryActivityCount, nil
	}
	return lotteryActivities, lotteryActivityCount, nil
}

func (ability *LotteryAbility) IncrementSalesVolume(ctx jet.Ctx, activityId uint, count int) error {
	return ability.lotteryActivityRepo.IncrementSalesVolume(ctx, activityId, count)
}

func (ability *LotteryAbility) ListActivityPrize(ctx jet.Ctx, params *api.PageParams) ([]*dto.LotteryActivityDTO, int64, error) {
	// 0. ability query
	lotteryActivities, lotteryActivityCount, err := ability.lotteryActivityRepo.List(params.Page, params.PageSize, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "ListActivityPrize error")
	}
	if lotteryActivities == nil || len(lotteryActivities) == 0 {
		return make([]*dto.LotteryActivityDTO, 0, 0), lotteryActivityCount, nil
	}
	// 1. 查询奖品信息
	lotteryActivityIds := stream.Of[*po.LotteryActivity, uint](lotteryActivities).
		Map(func(ele *po.LotteryActivity) uint { return ele.ID }).
		CollectToSlice()

	activityPrizeRelations, err := ability.relationRepo.Find("activity_id in (?)", lotteryActivityIds)
	if err != nil {
		ctx.Logger().Errorf("[LotteryAbility#ListActivity] find activity_prize_relation error, %v", err)
		return nil, 0, errors.Wrap(err, "ListActivity relationRepo error")
	}

	lotteryPrizeIds := stream.Of[*po.ActivityPrizeRelation, uint](activityPrizeRelations).
		Map(func(ele *po.ActivityPrizeRelation) uint { return ele.PrizeID }).
		CollectToSlice()

	lotteryPrizes, err := ability.lotteryPrizeRepo.Find("id in (?)", lotteryPrizeIds)
	if err != nil {
		ctx.Logger().Errorf("[LotteryAbility#ListActivity] find lottery_prize error, %v", err)
		return nil, 0, errors.Wrap(err, "ListActivity lotteryPrizeRepo error")
	}

	lotteryPrizeMap := utils.SliceToSingleMap[*po.LotteryPrize, uint](lotteryPrizes,
		func(ele *po.LotteryPrize) uint { return ele.ID })

	activityId2PrizeFunc := func(activityId uint) []*po.LotteryPrize {
		filterRelations := utils.Filter(activityPrizeRelations,
			func(ele *po.ActivityPrizeRelation) bool { return ele.ActivityID == activityId })
		prizes := stream.Of[*po.ActivityPrizeRelation, *po.LotteryPrize](filterRelations).
			Map(func(ele *po.ActivityPrizeRelation) *po.LotteryPrize {
				return lotteryPrizeMap.GetOrDefault(ele.PrizeID, nil)
			}).
			CollectToSlice()
		return utils.Filter(prizes, func(ele *po.LotteryPrize) bool { return ele != nil })
	}
	return stream.Of[*po.LotteryActivity, *dto.LotteryActivityDTO](lotteryActivities).
		Map(func(ele *po.LotteryActivity) *dto.LotteryActivityDTO {
			return &dto.LotteryActivityDTO{
				LotteryActivity: ele,
				LotteryPrizes:   activityId2PrizeFunc(ele.ID),
			}
		}).
		CollectToSlice(), int64(len(lotteryActivities)), nil
}

func (ability *LotteryAbility) DelActivity(ctx jet.Ctx, lotteryActivityId uint) error {
	// 0. 删除活动
	if err := ability.lotteryActivityRepo.RemoveByID(lotteryActivityId); err != nil {
		return errors.Wrap(err, "删除活动失败")
	}
	// 1. 查找关系
	relations, err := ability.relationRepo.FindByActivityId(ctx, lotteryActivityId)
	if err != nil {
		return errors.Wrap(err, "查找关系失败")
	}
	if relations == nil || len(relations) == 0 {
		return nil
	}
	// 删除奖品
	prizeIds := utils.Map(relations, func(relation *po.ActivityPrizeRelation) uint { return relation.PrizeID })

	count, err := ability.lotteryPrizeRepo.RemoveByPrizeIds(ctx, prizeIds)
	if err != nil {
		return errors.Wrap(err, "删除奖品失败")
	}
	ctx.Logger().Infof("[LotteryAbility#DelActivity] delete %d prizes", count)

	// 3. 删除关系
	if err = ability.relationRepo.DelByActivityId(ctx, lotteryActivityId); err != nil {
		return errors.Wrap(err, "删除关系失败")
	}

	return nil
}

func (ability *LotteryAbility) FindActivityPrizeByActivityId(ctx jet.Ctx, activityId uint) (*dto.LotteryActivityDTO, error) {
	// 0. ability query
	lotteryActivity, err := ability.lotteryActivityRepo.FindByID(activityId)
	if err != nil || lotteryActivity == nil || lotteryActivity.ID <= 0 {
		return nil, errors.Wrap(err, "ListActivityPrize error")
	}
	// 1. 查询奖品信息
	activityPrizeRelations, err := ability.relationRepo.Find("activity_id in (?)", lotteryActivity.ID)
	if err != nil {
		ctx.Logger().Errorf("[LotteryAbility#ListActivity] find activity_prize_relation error, %v", err)
		return nil, errors.Wrap(err, "ListActivity relationRepo error")
	}

	lotteryPrizeIds := stream.Of[*po.ActivityPrizeRelation, uint](activityPrizeRelations).
		Map(func(ele *po.ActivityPrizeRelation) uint { return ele.PrizeID }).
		CollectToSlice()

	lotteryPrizes, err := ability.lotteryPrizeRepo.Find("id in (?)", lotteryPrizeIds)
	if err != nil {
		ctx.Logger().Errorf("[LotteryAbility#ListActivity] find lottery_prize error, %v", err)
		return nil, errors.Wrap(err, "ListActivity lotteryPrizeRepo error")
	}
	slices.SortFunc(lotteryPrizes, func(a, b *po.LotteryPrize) int {
		return a.SortOrder - b.SortOrder
	})
	return &dto.LotteryActivityDTO{
		LotteryActivity: lotteryActivity,
		LotteryPrizes:   lotteryPrizes,
	}, nil
}

func (ability *LotteryAbility) AddLotteryRecords(ctx jet.Ctx, lotteryRecords *po.LotteryRecords) error {
	err := ability.lotteryRecordsRepo.InsertOne(lotteryRecords)
	if err != nil {
		return errors.Wrap(err, "insert lottery records error")
	}
	return nil
}

func (ability *LotteryAbility) FindFallbackPrize(ctx jet.Ctx, activityId uint) (*po.LotteryPrize, error) {
	activityPO, err := ability.lotteryActivityRepo.FindByID(activityId)
	if err != nil || activityPO == nil {
		ctx.Logger().Errorf(
			"[LotteryAbility#FindFallbackPrize] find activityId:%v activity_prize_relation error, %v", activityId, err)
		return nil, errors.Wrap(err, "find fallback prize error")
	}
	fallbackPrizeId := activityPO.FallbackPrizeId
	prizePO, err := ability.lotteryPrizeRepo.FindByID(fallbackPrizeId)
	if err != nil {
		ctx.Logger().Errorf(
			"[LotteryAbility#FindFallbackPrize] find activityId:%v activity_prize_relation error, %v", activityId, err)
		return nil, errors.Wrap(err, "find fallback prize error")
	}
	return prizePO, nil
}
