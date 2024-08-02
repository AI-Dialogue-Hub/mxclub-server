package dto

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
