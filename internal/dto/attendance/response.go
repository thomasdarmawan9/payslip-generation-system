package attendance

type SubmitAttendanceResponse struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Date   string `json:"date"`
	Status string `json:"status"` // "created" atau "already_exists"
}
