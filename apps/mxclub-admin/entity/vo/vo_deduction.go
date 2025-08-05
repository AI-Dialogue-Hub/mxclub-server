package vo

import "time"

type DeductionVO struct {
	ID              uint      `json:"id"`
	DasherId        uint      `json:"dasher_id"`
	UserID          uint      `json:"user_id"`
	UserInfo        string    `json:"user_info"`
	ConfirmPersonId uint      `json:"confirm_person_id"`
	Amount          float64   `json:"amount"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
