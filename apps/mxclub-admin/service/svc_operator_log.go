package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/dto"
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
	operatorPO := &po.OperatorLogPO{
		Type:     operatorDTO.Type,
		Remarks:  operatorDTO.Remarks,
		UserId:   0,
		UserName: "",
	}
	logger.Infof("[OperatorLogService#DoLog] handle operator log: %v", utils.ObjToJsonStr(operatorPO))
	if err := svc.operatorLogRepo.InsertOne(operatorPO); err != nil {
		logger.Errorf("[OperatorLogService#DoLog] InsertOne ERROR, %v", err)
	}
}
