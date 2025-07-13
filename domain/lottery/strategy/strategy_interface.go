package strategy

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/lottery/entity/dto"
)

// ILotteryStrategy 主策略接口
type ILotteryStrategy interface {
	DoDraw(ctx jet.Ctx, drawInfo *dto.LotteryStrategyDrawDTO) (*dto.LotteryStrategyDrawResultDTO, error)
}
