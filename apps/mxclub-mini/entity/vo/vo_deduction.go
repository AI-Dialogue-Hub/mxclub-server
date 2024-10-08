package vo

import "time"

type DeductionVO struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	ConfirmPersonId uint      `json:"confirm_person_id"`
	Amount          float64   `json:"executor_price"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"purchase_date"`
}
