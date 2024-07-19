package vo

type AssistantApplicationVO struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Phone        string `json:"phone"`
	MemberNumber int64  `json:"member_number"`
	Name         string `json:"name"`
	Status       string `json:"status"`
}
