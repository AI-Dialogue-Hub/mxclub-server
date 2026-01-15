package strategy

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"mxclub/domain/lottery/ability"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/utils"
)

var _ ILotteryStrategyHook = (*LotteryStrategyHook)(nil)

type LotteryStrategyHook struct {
}

func NewLotteryStrategyHook(lotteryAbility ability.ILotteryAbility) *LotteryStrategyHook {
	return &LotteryStrategyHook{}
}

const (
	// 三等奖连续抽中次数上线，达到默认值，执行抽奖规则
	lotteryPrizeLevelThirdQuota = 3
	// lottery:miss_count:{user_id}:{activity_id}:{prize_level}
	missKeyTemplate = "lottery:miss_count:%v:%v:%v"
)

// BeforeDraw 在抽奖前执行，嵌入抽奖规则逻辑
//
// 规则：
//   - 连续抽中三次三等奖后，下一次必中二等奖
//   - 中二等奖后，计数器清零
//
// 参数：
//   - ctx: 上下文，用于传递请求信息
//
// 返回值：
//   - error: 如果规则检查失败，返回错误
func (l *LotteryStrategyHook) BeforeDraw(
	ctx jet.Ctx, beforeDrawDTO *dto.LotteryStrategyBeforeDrawDTO,
) (*dto.LotteryStrategyDrawResultDTO, error) {
	defer utils.RecoverWithPrefix(ctx, "LotteryStrategyHook#BeforeDraw")
	// 0. 获取用户连抽不中记录
	missVal, err := fetchMissVal(ctx, beforeDrawDTO.UserId, beforeDrawDTO.ActivityId)
	if err != nil {
		return nil, err
	}
	// 1. 执行策略
	// 1.1 连续抽中三次三等奖后，中一次二等奖
	if missVal >= lotteryPrizeLevelThirdQuota {
		lotteryAbility := ability.FetchLotteryAbilityInstance()
		// 1.1.1 获取中奖奖品
		prize, err := lotteryAbility.FindFallbackPrize(ctx, beforeDrawDTO.ActivityId)
		if err != nil {
			ctx.Logger().Errorf("get lottery records error,err info:%v", err)
			return nil, errors.Wrap(err, "get lottery records error")
		}
		// 1.1.2 落库
		err = lotteryAbility.AddLotteryRecords(ctx, &po.LotteryRecords{
			ActivityId: beforeDrawDTO.ActivityId,
			PrizeId:    prize.ID,
			UserId:     beforeDrawDTO.UserId,
			OrderId:    "",
			ActivityPrizeSnapshot: utils.ObjToJsonStr(
				struct {
					ActivityId uint             `json:"activity_id"`
					Prize      *po.LotteryPrize `json:"prize"`
				}{ActivityId: beforeDrawDTO.ActivityId, Prize: prize},
			),
		})
		if err != nil {
			ctx.Logger().Errorf("insert lottery records error,err info:%v", err)
			return nil, errors.Wrap(err, "add lottery records error")
		}

		return &dto.LotteryStrategyDrawResultDTO{
			PrizeId:      prize.ID,
			LotteryPrize: prize,
		}, nil
	}
	return nil, nil
}

func fetchMissVal(
	ctx jet.Ctx,
	userId uint,
	activityId uint,
) (int, error) {
	missKey := fmt.Sprintf(missKeyTemplate, userId, activityId, enum.PrizeLevelThird)
	val, err := xredis.GetInt(missKey)
	if err != nil || val <= 0 {
		ctx.Logger().Errorf("get miss count error,missKey => %v, err info:%v", missKey, err)
		return 0, nil
	}
	return val, nil
}

func (l *LotteryStrategyHook) AfterDraw(ctx jet.Ctx, afterDrawDTO *dto.LotteryStrategyAfterDrawDTO) {
	defer utils.RecoverWithPrefix(ctx, "LotteryStrategyHook#AfterDraw")
	ctx.Logger().Infof("do after draw, afterDrawDTO:%v", utils.ObjToJsonStr(afterDrawDTO))

	lotteryAbility := ability.FetchLotteryAbilityInstance()
	missKey := fmt.Sprintf(missKeyTemplate, afterDrawDTO.UserId, afterDrawDTO.ActivityId, enum.PrizeLevelThird)

	// 根据中奖等级处理计数器
	switch afterDrawDTO.PrizeLevel {
	case enum.PrizeLevelThird:
		// 中三等奖：增加计数
		missVal, err := fetchMissVal(ctx, afterDrawDTO.UserId, afterDrawDTO.ActivityId)
		if err != nil {
			ctx.Logger().Errorf("get miss count error,err info:%v", err)
			return
		}
		if missVal < lotteryPrizeLevelThirdQuota {
			// 增加值
			incrVal, err := xredis.Incr(missKey)
			ctx.Logger().Infof("useId:%v, activityId:%v, incr count: %v, missKey => %v, err info:%v",
				afterDrawDTO.UserId, afterDrawDTO.ActivityId, incrVal, missKey, err)
		}
	case enum.PrizeLevelSecond, enum.PrizeLevelFirst:
		// 中一等奖或二等奖：清零计数器
		if err := xredis.Del(missKey); err != nil {
			ctx.Logger().Warnf("failed to clear miss count, missKey => %v, err info:%v", missKey, err)
		} else {
			ctx.Logger().Infof("cleared miss count after winning %v, userId:%v, activityId:%v",
				afterDrawDTO.PrizeLevel, afterDrawDTO.UserId, afterDrawDTO.ActivityId)
		}
	}

	// 1. 奖品中奖落库
	err := lotteryAbility.AddLotteryRecords(ctx, &po.LotteryRecords{
		ActivityId:            afterDrawDTO.ActivityId,
		PrizeId:               afterDrawDTO.PrizeId,
		UserId:                afterDrawDTO.UserId,
		OrderId:               "",
		ActivityPrizeSnapshot: afterDrawDTO.ActivityPrizeSnapshot,
	})
	if err != nil {
		ctx.Logger().Errorf("insert lottery records error,err info:%v", err)
		return
	}
	// 2. 销量增加
	err = lotteryAbility.IncrementSalesVolume(ctx, afterDrawDTO.ActivityId, 1)
	if err != nil {
		ctx.Logger().Errorf("increment sales volume error,err info:%v", err)
	}
}
