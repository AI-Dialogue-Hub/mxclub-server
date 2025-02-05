package enum

type TransferEnum int

const (
	Transfer_PENDING TransferEnum = iota
	Transfer_SUCCESS
	Transfer_REJECT
)

func (t TransferEnum) IsValid() bool {
	return t <= Transfer_REJECT && t >= Transfer_PENDING
}
