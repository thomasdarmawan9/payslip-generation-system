package attendance

type SubmitAttendanceRequest struct {
	// Optional, default = today (WIB). Format YYYY-MM-DD
	Date string `json:"date" binding:"omitempty,datetime=2006-01-02"`
}
