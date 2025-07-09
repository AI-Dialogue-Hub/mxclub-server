package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/lottery/ability"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewLotteryService)
}

type LotteryService struct {
	lotteryPrizeRepo repo.ILotteryPrizeRepo
	lotteryActivity  ability.ILotteryAbility
}

func NewLotteryService(lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivity ability.ILotteryAbility) *LotteryService {
	return &LotteryService{lotteryPrizeRepo: lotteryPrizeRepo, lotteryActivity: lotteryActivity}
}

func (svc *LotteryService) ListLotteryPrize(ctx jet.Ctx, params *api.PageParams) ([]*vo.LotteryActivityPrizeVO, int64, error) {
	listActivity, count, err := svc.lotteryActivity.ListActivityPrize(ctx, params)
	if err != nil {
		ctx.Logger().Errorf("[LotteryService#ListLotteryPrize] ERROR:%v", err)
		return nil, 0, errors.New("活动获取错误")
	}
	lotteryActivityPrizeVOS := utils.Map(listActivity, func(activityDTO *dto.LotteryActivityDTO) *vo.LotteryActivityPrizeVO {
		return &vo.LotteryActivityPrizeVO{
			LotteryActivity: utils.MustCopy[vo.LotteryActivityVO](activityDTO.LotteryActivity),
			LotteryPrizes:   utils.CopySlice[*po.LotteryPrize, *vo.LotteryPrizeVO](activityDTO.LotteryPrizes),
		}
	})
	return lotteryActivityPrizeVOS, count, nil
}

func (svc *LotteryService) FindActivityPrizeByActivityId(ctx jet.Ctx, activityId int) (*vo.LotteryActivityPrizeVO, error) {
	activityDTO, err := svc.lotteryActivity.FindActivityPrizeByActivityId(ctx, uint(activityId))
	if err != nil {
		return nil, errors.New("活动获取错误")
	}
	return utils.MustCopy[vo.LotteryActivityPrizeVO](activityDTO), nil
}
