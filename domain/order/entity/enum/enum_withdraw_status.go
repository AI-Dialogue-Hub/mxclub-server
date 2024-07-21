package enum

type WithdrawalStatus string

func Initiated() string {
	return "initiated"
}

func Completed() string {
	return "completed"
}

func Reject() string {
	return "reject"
}
