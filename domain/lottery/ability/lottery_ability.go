package ability

import (
	"github.com/fengyuan-liang/GoKit/collection/stream"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewLotteryActivity)
}

// ILotteryAbility 是领域根，不是仓储实现
type ILotteryAbility interface {
	FetchLotteryPrizeType() *dto.LotteryPrizeTypeDTO
	AddPrize(ctx jet.Ctx, activityId uint, lotteryPrize *po.LotteryPrize) error
	DelPrize(ctx jet.Ctx, lotteryPrizeId uint) error
	AddOrUpdateActivity(ctx jet.Ctx, lotteryability *po.LotteryActivity) error
	ListActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error)
	ListActivityPrize(ctx jet.Ctx, params *api.PageParams) ([]*dto.LotteryActivityDTO, int64, error)
	DelActivity(ctx jet.Ctx, lotteryActivityId uint) error
	FindActivityPrizeByActivityId(ctx jet.Ctx, activityId uint) (*dto.LotteryActivityDTO, error)
}

type LotteryAbility struct {
	lotteryRepo         repo.ILotteryRepo
	lotteryPrizeRepo    repo.ILotteryPrizeRepo
	lotteryActivityRepo repo.ILotteryActivityRepo
	relationRepo        repo.ILotteryActivityPrizeRelationRepo
}

func NewLotteryActivity(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivityRepo repo.ILotteryActivityRepo,
	lotteryRepo repo.ILotteryRepo,
	relationRepo repo.ILotteryActivityPrizeRelationRepo) ILotteryAbility {
	return &LotteryAbility{
		lotteryRepo:         lotteryRepo,
		lotteryPrizeRepo:    lotteryPrizeRepo,
		lotteryActivityRepo: lotteryActivityRepo,
		relationRepo:        relationRepo,
	}
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
		if err != gorm.ErrRecordNotFound {
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
	return &dto.LotteryActivityDTO{
		LotteryActivity: lotteryActivity,
		LotteryPrizes:   lotteryPrizes,
	}, nil
}
