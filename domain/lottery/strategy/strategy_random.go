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

const RandomStrategy = "randomLotteryStrategy"

func init() {
	jet.Provide(func(lotteryAbility ability.ILotteryAbility) {
		utils.IfNotNilPanic(
			OfferLotteryStrategy(
				RandomStrategy,
				NewLotteryStrategyRandom(
					NewLotteryStrategyHook(lotteryAbility),
					lotteryAbility,
				),
			),
		)
	})
}

// LotteryStrategyRandom 具体策略
type LotteryStrategyRandom struct {
	*LotteryStrategyBase // 组合基础模板
	lotteryAbility       ability.ILotteryAbility
}

// NewLotteryStrategyRandom 构造函数（注入钩子）
func NewLotteryStrategyRandom(
	hook ILotteryStrategyHook, lotteryAbility ability.ILotteryAbility,
) *LotteryStrategyRandom {
	return &LotteryStrategyRandom{
		lotteryAbility:      lotteryAbility,
		LotteryStrategyBase: &LotteryStrategyBase{LotteryStrategyHook: hook},
	}
}

// handleDraw 实现具体逻辑
func (l *LotteryStrategyRandom) handleDraw(
	ctx jet.Ctx,
	drawInfo *dto.LotteryStrategyDrawDTO,
) (*dto.LotteryStrategyDrawResultDTO, error) {
	// 0. 获取活动信息
	activityDTO, err := l.lotteryAbility.FindActivityPrizeByActivityId(ctx, drawInfo.ActivityId)
	if err != nil {
		ctx.Logger().Errorf("FindActivityPrizeByActivityId error,err info:%v", err)
		return nil, errors.Wrap(err, "FindActivityPrizeByActivityId error")
	}
	drawInfo.ActivityPrizeSnapshot = utils.ObjToJsonStr(activityDTO)
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
		ActivityId:                   drawInfo.ActivityId,
		IsIncreaseLotteryProbability: true,
		LotteryPrize:                 finalPrize,
		PrizeId:                      finalPrize.ID,
		UserId:                       drawInfo.UserId,
	}, nil
}
