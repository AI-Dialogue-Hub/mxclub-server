package activity

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewLotteryActivity)
}

// ILotteryActivity 是领域根，不是仓储实现
type ILotteryActivity interface {
	FetchLotteryPrizeType() *dto.LotteryPrizeTypeDTO
	AddPrize(ctx jet.Ctx, activityId uint, po *po.LotteryPrize) error
	DelPrize(ctx jet.Ctx, lotteryPrizeId uint) error
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
		return err
	}
	return nil
}
