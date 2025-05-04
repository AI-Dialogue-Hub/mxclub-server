package service

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/dto"
	"mxclub/apps/mxclub-admin/middleware"
	"mxclub/domain/operator/po"
	"mxclub/domain/operator/repo"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewOperatorLogService)
}

type OperatorLogService struct {
	operatorLogRepo repo.IOperatorRepo
}

func NewOperatorLogService(operatorLogRepo repo.IOperatorRepo) *OperatorLogService {
	return &OperatorLogService{operatorLogRepo: operatorLogRepo}
}

func (svc OperatorLogService) DoLog(ctx jet.Ctx, operatorDTO *dto.OperatorLogDTO) {
	var (
		logger = ctx.Logger()
	)
	userInfo, err := middleware.FetchUserInfoByCtx(ctx.FastHttpCtx())
	if err != nil {
		logger.Errorf("[OperatorLogService#DoLog] failed FetchUserInfoByCtx, %v", err)
		return
	}
	operatorPO := &po.OperatorLogPO{
		Type:     operatorDTO.Type,
		Remarks:  operatorDTO.Remarks,
		UserId:   userInfo.ID,
		UserName: userInfo.Name,
	}
	logger.Infof("[OperatorLogService#DoLog] handle operator log: %v", utils.ObjToJsonStr(operatorPO))
	if err = svc.operatorLogRepo.InsertOne(operatorPO); err != nil {
		logger.Errorf("[OperatorLogService#DoLog] InsertOne ERROR, %v", err)
	}
}

func (svc OperatorLogService) FindRefundOrderLog(ctx jet.Ctx, orderId string) (*dto.OperatorLogDTO, error) {
	operatorLogPO, err := svc.operatorLogRepo.FindByBizId(orderId)
	if err != nil {
		ctx.Logger().Errorf("FindRefundOrderLog ERROR, %v", err)
		return nil, fmt.Errorf("cannot find by orderId:%v", orderId)
	}
	return utils.MustCopy[dto.OperatorLogDTO](operatorLogPO), nil
}
