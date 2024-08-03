package dto

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/api"
)

type DeductionDTO struct {
	*api.PageParams `json:"page_params"`
	Ge              string `json:"ge"` // GREATER THAN大于
	Le              string `json:"le"` // LESS THAN小于
	Status          *enum.DeductStatus
}
