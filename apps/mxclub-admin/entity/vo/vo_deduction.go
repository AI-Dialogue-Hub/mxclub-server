package vo

import "mxclub/domain/order/entity/enum"

type DeductionVO struct {
	ID              uint              `json:"id"`
	UserID          uint              `json:"user_id"`
	ConfirmPersonId uint              `json:"confirm_person_id"`
	Amount          float64           `json:"amount"`
	Reason          string            `json:"reason"`
	Status          enum.DeductStatus `json:"status"`
}
