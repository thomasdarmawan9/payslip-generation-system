package overtime

type SubmitOvertimeRequest struct {
	// default = today (WIB), format YYYY-MM-DD
	Date  string  `json:"date"  binding:"omitempty,datetime=2006-01-02"`
	Hours float64 `json:"hours" binding:"required,gt=0,lte=3"`
}
