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
