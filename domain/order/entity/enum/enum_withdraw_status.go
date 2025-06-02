package enum

type WithdrawalStatus string

func Initiated() WithdrawalStatus {
	return "initiated"
}

func Completed() string {
	return "completed"
}

func Reject() string {
	return "reject"
}

func (r WithdrawalStatus) DisplayName() string {
	if r == "initiated" {
		return "申请中"
	}
	if r == "completed" {
		return "已提现"
	}
	if r == "reject" {
		return "已拒绝"
	}
	return "订单状态未知"
}
