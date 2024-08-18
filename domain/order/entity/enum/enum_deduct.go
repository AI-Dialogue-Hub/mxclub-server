package enum

type DeductStatus string

const (
	Deduct_PENDING DeductStatus = "PENDING"
	Deduct_SUCCESS DeductStatus = "SUCCESS"
	Deduct_REJECT  DeductStatus = "REJECT"
)

var m = map[DeductStatus]string{
	Deduct_PENDING: "待处罚",
	Deduct_SUCCESS: "已处罚",
	Deduct_REJECT:  "已拒绝",
}

var reverseMap = func() map[string]DeductStatus {
	tempMap := make(map[string]DeductStatus)
	for k, v := range m {
		tempMap[v] = k
	}
	return tempMap
}()

func (d DeductStatus) DisPlayName() string {
	return m[d]
}

func Of(str string) DeductStatus {
	return reverseMap[str]
}
