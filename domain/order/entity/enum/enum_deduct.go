package enum

type DeductStatus string

const (
	Deduct_PENDING DeductStatus = "PENDING"
	Deduct_SUCCESS DeductStatus = "SUCCESS"
)

var m = map[DeductStatus]string{
	Deduct_PENDING: "待处罚",
	Deduct_SUCCESS: "已处罚",
}

func (d DeductStatus) DisPlayName() string {
	return m[d]
}
