package activity

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

// ILotteryActivity 是领域根，不是仓储实现
type ILotteryActivity interface {
	FetchLotteryPrizeType() *dto.LotteryPrizeTypeDTO
	AddPrize(ctx jet.Ctx, activityId uint, lotteryPrize *po.LotteryPrize) error
	DelPrize(ctx jet.Ctx, lotteryPrizeId uint) error
	AddOrUpdateActivity(ctx jet.Ctx, lotteryActivity *po.LotteryActivity) error
	ListActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error)
	ListActivityPrize(ctx jet.Ctx, params *api.PageParams) ([]*dto.LotteryActivityDTO, int64, error)
	DelActivity(ctx jet.Ctx, lotteryActivityId uint) error
}

type LotteryActivity struct {
	lotteryRepo         repo.ILotteryRepo
	lotteryPrizeRepo    repo.ILotteryPrizeRepo
	lotteryActivityRepo repo.ILotteryActivityRepo
	relationRepo        repo.ILotteryActivityPrizeRelationRepo
}

func NewLotteryActivity(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivityRepo repo.ILotteryActivityRepo,
	lotteryRepo repo.ILotteryRepo,
	relationRepo repo.ILotteryActivityPrizeRelationRepo) ILotteryActivity {
	return &LotteryActivity{
		lotteryRepo:         lotteryRepo,
		lotteryPrizeRepo:    lotteryPrizeRepo,
		lotteryActivityRepo: lotteryActivityRepo,
		relationRepo:        relationRepo,
	}
}

func (activity *LotteryActivity) FetchLotteryPrizeType() *dto.LotteryPrizeTypeDTO {
	return &dto.LotteryPrizeTypeDTO{PrizeType: enum.PrizeTypeNames}
}

func (activity *LotteryActivity) AddPrize(ctx jet.Ctx, activityId uint, lotteryPrize *po.LotteryPrize) error {
	// 1. 插入奖品
	if err := activity.lotteryPrizeRepo.InsertOne(lotteryPrize); err != nil {
		return err
	}
	ctx.Logger().Infof("[LotteryActivity#AddPrize] success, lotteryPrize=%v", utils.ObjToJsonStr(lotteryPrize))
	// 2. 插入活动-奖品关系
	if err := activity.relationRepo.InsertOne(&po.ActivityPrizeRelation{
		ActivityID: activityId,
		PrizeID:    lotteryPrize.ID,
	}); err != nil {
		return err
	}
	return nil
}

func (activity *LotteryActivity) DelPrize(ctx jet.Ctx, lotteryPrizeId uint) error {
	// 1. 删除奖品
	if err := activity.lotteryPrizeRepo.RemoveByID(lotteryPrizeId); err != nil {
		ctx.Logger().Errorf("[LotteryActivity#DelPrize] err:%v", err)
		return err
	}
	// 2. 删除跟活动的关联关系
	if err := activity.relationRepo.DelByPrizeId(ctx, lotteryPrizeId); err != nil {
		ctx.Logger().Errorf("[LotteryActivity#DelPrize] err:%v", err)
		if err != gorm.ErrRecordNotFound {
			return nil
		}
	}
	return nil
}

func (activity *LotteryActivity) AddOrUpdateActivity(ctx jet.Ctx, lotteryActivity *po.LotteryActivity) error {
	if lotteryActivity.ID > 0 {
		if err := activity.lotteryActivityRepo.UpdateById(lotteryActivity, lotteryActivity.ID); err != nil {
			ctx.Logger().Errorf("[LotteryActivity#AddOrUpdateActivity] err:%v", err)
			return err
		}
		return nil
	}
	return activity.lotteryActivityRepo.InsertOne(lotteryActivity)
}

func (activity *LotteryActivity) ListActivity(ctx jet.Ctx, params *api.PageParams) ([]*po.LotteryActivity, int64, error) {
	lotteryActivities, lotteryActivityCount, err := activity.lotteryActivityRepo.List(params.Page, params.PageSize, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "ListActivity error")
	}
	if lotteryActivities == nil || len(lotteryActivities) == 0 {
		return make([]*po.LotteryActivity, 0, 0), lotteryActivityCount, nil
	}
	return lotteryActivities, lotteryActivityCount, nil
}

func (activity *LotteryActivity) ListActivityPrize(ctx jet.Ctx, params *api.PageParams) ([]*dto.LotteryActivityDTO, int64, error) {
	// 0. activity query
	lotteryActivities, lotteryActivityCount, err := activity.lotteryActivityRepo.List(params.Page, params.PageSize, nil)
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

	activityPrizeRelations, err := activity.relationRepo.Find("activity_id in (?)", lotteryActivityIds)
	if err != nil {
		ctx.Logger().Errorf("[LotteryActivity#ListActivity] find activity_prize_relation error, %v", err)
		return nil, 0, errors.Wrap(err, "ListActivity relationRepo error")
	}

	lotteryPrizeIds := stream.Of[*po.ActivityPrizeRelation, uint](activityPrizeRelations).
		Map(func(ele *po.ActivityPrizeRelation) uint { return ele.PrizeID }).
		CollectToSlice()

	lotteryPrizes, err := activity.lotteryPrizeRepo.Find("id in (?)", lotteryPrizeIds)
	if err != nil {
		ctx.Logger().Errorf("[LotteryActivity#ListActivity] find lottery_prize error, %v", err)
		return nil, 0, errors.Wrap(err, "ListActivity lotteryPrizeRepo error")
	}

	lotteryPrizeMap := utils.SliceToSingleMap[*po.LotteryPrize, uint](lotteryPrizes,
		func(ele *po.LotteryPrize) uint { return ele.ID })

	activityId2PrizeFunc := func(activityId uint) []*po.LotteryPrize {
		filterRelations := utils.Filter(activityPrizeRelations,
			func(ele *po.ActivityPrizeRelation) bool { return ele.ID == activityId })
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

func (activity *LotteryActivity) DelActivity(ctx jet.Ctx, lotteryActivityId uint) error {
	// 0. 删除活动
	if err := activity.lotteryActivityRepo.RemoveByID(lotteryActivityId); err != nil {
		return errors.Wrap(err, "删除活动失败")
	}
	// 1. 查找关系
	relations, err := activity.relationRepo.FindByActivityId(ctx, lotteryActivityId)
	if err != nil {
		return errors.Wrap(err, "查找关系失败")
	}
	if relations == nil || len(relations) == 0 {
		return nil
	}
	// 删除奖品
	prizeIds := utils.Map(relations, func(relation *po.ActivityPrizeRelation) uint { return relation.PrizeID })

	count, err := activity.lotteryPrizeRepo.RemoveByPrizeIds(ctx, prizeIds)
	if err != nil {
		return errors.Wrap(err, "删除奖品失败")
	}
	ctx.Logger().Infof("[LotteryActivity#DelActivity] delete %d prizes", count)

	// 3. 删除关系
	if err = activity.relationRepo.DelByActivityId(ctx, lotteryActivityId); err != nil {
		return errors.Wrap(err, "删除关系失败")
	}

	return nil
}
