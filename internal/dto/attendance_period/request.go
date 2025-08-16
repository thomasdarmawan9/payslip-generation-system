package attendance_period

type CreatePeriodRequest struct {
	Name      string `json:"name"`                 // optional
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date"   binding:"required,datetime=2006-01-02"`
}
