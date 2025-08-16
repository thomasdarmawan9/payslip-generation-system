package attendance_period

type PeriodResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`
}
