package strategy

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"mxclub/domain/lottery/entity/dto"
)

// ILotteryStrategyHook 定义钩子接口（公开所有方法）
type ILotteryStrategyHook interface {
	BeforeDraw(ctx jet.Ctx, beforeDrawDTO *dto.LotteryStrategyBeforeDrawDTO) (*dto.LotteryStrategyDrawResultDTO, error)
	AfterDraw(ctx jet.Ctx, afterDrawDTO *dto.LotteryStrategyAfterDrawDTO)
}

// LotteryStrategyBase 基础模板（组合钩子接口）
type LotteryStrategyBase struct {
	ILotteryStrategyHook                      // 组合钩子接口
	LotteryStrategyHook  ILotteryStrategyHook // 保存子类钩子接口
}

// DoDraw 模板方法（固定流程）
func (s *LotteryStrategyBase) DoDraw(
	ctx jet.Ctx,
	drawInfo *dto.LotteryStrategyDrawDTO,
) (*dto.LotteryStrategyDrawResultDTO, error) {
	// 0. BeforeDraw
	beforeDrawDTO := drawInfo.WrapBeforeDrawDTO()
	if result, err := s.BeforeDraw(ctx, beforeDrawDTO); err == nil && result != nil {
		return result, nil
	}
	// 1. 调用子类实现的 handleDraw
	result, err := s.handleDraw(ctx, drawInfo)
	if err != nil {
		ctx.Logger().Errorf("handleDraw error: %v", err)
		return nil, errors.Wrap(err, "handleDraw error")
	}

	// 2. AfterDraw
	s.AfterDraw(ctx, drawInfo.WrapAfterDrawDTO())
	return result, nil
}

// handleDraw 内部方法（由子类实现）
//
// @see LotteryStrategyRandom
func (s *LotteryStrategyBase) handleDraw(
	ctx jet.Ctx, drawInfo *dto.LotteryStrategyDrawDTO,
) (*dto.LotteryStrategyDrawResultDTO, error) {
	panic("implement me") // 强制子类实现
}
