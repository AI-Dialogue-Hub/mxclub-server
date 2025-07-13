package strategy

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"math/rand"
	"mxclub/domain/lottery/ability"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/utils"
)

// ILotteryStrategyHook 定义钩子接口（公开所有方法）
type ILotteryStrategyHook interface {
	BeforeDraw(ctx jet.Ctx, beforeDrawDTO *dto.LotteryStrategyBeforeDrawDTO) (*dto.LotteryStrategyDrawResultDTO, error)
	AfterDraw(ctx jet.Ctx, afterDrawDTO *dto.LotteryStrategyAfterDrawDTO)
}

// LotteryStrategyBase 基础模板（组合钩子接口）
type LotteryStrategyBase struct {
	hook ILotteryStrategyHook // 保存子类钩子接口
}

// DoDraw 模板方法（固定流程）
func (s *LotteryStrategyBase) DoDraw(
	ctx jet.Ctx,
	drawInfo *dto.LotteryStrategyDrawDTO,
) (*dto.LotteryStrategyDrawResultDTO, error) {
	defer utils.RecoverAndLogError(ctx)
	// 0. BeforeDraw
	beforeDrawDTO := drawInfo.WrapBeforeDrawDTO()
	if result, err := s.hook.BeforeDraw(ctx, beforeDrawDTO); err == nil && result != nil {
		return result, nil
	}
	// 1. 调用子类实现的 handleInnerDraw
	result, err := s.handleDraw(ctx, drawInfo)
	if err != nil {
		ctx.Logger().Errorf("handleInnerDraw error: %v", err)
		return nil, errors.Wrap(err, "handleInnerDraw error")
	}

	// 2. AfterDraw
	s.hook.AfterDraw(ctx, &dto.LotteryStrategyAfterDrawDTO{
		UserId:                drawInfo.UserId,
		ActivityId:            drawInfo.ActivityId,
		PrizeId:               result.LotteryPrize.ID,
		PrizeLevel:            result.LotteryPrize.PrizeLevel,
		ActivityPrizeSnapshot: result.ActivityPrizeSnapshot,
	})
	return result, nil
}

// handleDraw 内部方法（由子类实现）
//
// @see LotteryStrategyRandom
func (s *LotteryStrategyBase) handleDraw(
	ctx jet.Ctx, drawInfo *dto.LotteryStrategyDrawDTO,
) (*dto.LotteryStrategyDrawResultDTO, error) {
	lotteryAbility := ability.FetchLotteryAbilityInstance()
	// 0. 获取活动信息
	activityDTO, err := lotteryAbility.FindActivityPrizeByActivityId(ctx, drawInfo.ActivityId)
	if err != nil {
		ctx.Logger().Errorf("FindActivityPrizeByActivityId error,err info:%v", err)
		return nil, errors.Wrap(err, "FindActivityPrizeByActivityId error")
	}
	prizes := activityDTO.LotteryPrizes
	if prizes == nil || len(prizes) <= 0 {
		ctx.Logger().Errorf("activity:%v, no prizes", drawInfo.ActivityId)
		return nil, errors.New("no prizes")
	}
	var (
		totalProb  = 0.0
		finalPrize *po.LotteryPrize
	)
	// 1. 过滤概率为0的奖品
	prizes = utils.Filter(prizes, func(prize *po.LotteryPrize) bool {
		totalProb += prize.ActualProbability
		return prize.ActualProbability > 0
	})
	totalProb = utils.RoundToTwoDecimalPlaces(totalProb / 100.0)
	// 2. 验证概率
	if totalProb <= 0 || totalProb > 1.0 {
		ctx.Logger().Errorf(
			"activity:%v, invalid probability distribution, totalProb is %v", drawInfo.ActivityId, totalProb)
		return nil, errors.New("invalid probability distribution")
	}

	// 3. 抽奖逻辑
	randVal := rand.Float64() * totalProb
	cumulativeProb := 0.0
	for _, prize := range prizes {
		cumulativeProb += prize.ActualProbability
		if randVal <= cumulativeProb {
			finalPrize = prize
			break
		}
	}
	// 4. 处理未中奖情况，理论上应该不会出现
	if finalPrize == nil {
		ctx.Logger().Errorf("activity:%v, no prizes", drawInfo.ActivityId)
		// 选择第一个三等奖作为fallback返回
		if findFirst, ok := utils.FindFirst(prizes, func(in *po.LotteryPrize) bool {
			return in.PrizeLevel == enum.PrizeLevelThird
		}); ok {
			finalPrize = findFirst
		}
		return nil, errors.New("no prizes")
	}
	return &dto.LotteryStrategyDrawResultDTO{
		IsIncreaseLotteryProbability: true,
		LotteryPrize:                 finalPrize,
		PrizeId:                      finalPrize.ID,
		ActivityPrizeSnapshot:        utils.ObjToJsonStr(activityDTO),
		PrizeIndex:                   finalPrize.SortOrder,
	}, nil
}
