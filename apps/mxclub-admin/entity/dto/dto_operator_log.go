package dto

import "mxclub/domain/operator/entity/enum"

type OperatorLogDTO struct {
	Type    enum.OperatorEnum
	BizId   string // 用来定位操作表具体行的数据，例如 某一行订单数据
	Remarks string // 操作记录
}

func NewOrderRemoveOperatorLogDTO(bizId, remarks string) *OperatorLogDTO {
	return &OperatorLogDTO{
		Type:    enum.OrderRemove,
		BizId:   bizId,
		Remarks: remarks,
	}
}

func NewOrderRefundsOperatorLogDTO(bizId, remarks string) *OperatorLogDTO {
	return &OperatorLogDTO{
		Type:    enum.OrderRefunds,
		BizId:   bizId,
		Remarks: remarks,
	}
}
