package dto

import (
	"mxclub/domain/order/entity/enum"
	"mxclub/pkg/api"
)

type OrderExecutorDTO struct {
	ExecutorType uint // Two or three
	ExecutorName string
	ExecutorId   uint
	OrderId      uint
}

type WithdrawAbleAmountDTO struct {
	DasherID int
}

type FinishOrderDTO struct {
	Id            uint
	Images        []string
	ExecutorNum   int
	ExecutorPrice float64
	CutRate       float64 // aka 0.2
}

type WithdrawListDTO struct {
	*api.PageParams `json:"page_params"`
	Ge              string `json:"ge"` // GREATER THAN大于
	Le              string `json:"le"` // LESS THAN小于
	Status          *enum.WithdrawalStatus
	UserId          uint
}

type ListByOrderStatusDTO struct {
	Status       enum.OrderStatus
	PageParams   *api.PageParams
	Ge           string
	Le           string
	MemberNumber int
	UserId       uint
	ExecutorName string
	IsDasher     bool
}
