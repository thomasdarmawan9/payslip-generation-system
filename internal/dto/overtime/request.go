package overtime

type SubmitOvertimeResponse struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Date   string `json:"date"`
	Hours  string `json:"hours"` // string biar rapi saat format (2 desimal)
	Status string `json:"status"` // created | already_exists
}