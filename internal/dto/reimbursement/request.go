package reimbursement

type CreateReimbursementRequest struct {
	// default = today (WIB), format YYYY-MM-DD
	Date        string  `json:"date"        binding:"omitempty,datetime=2006-01-02"`
	Amount      float64 `json:"amount"      binding:"required,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=255"`
}
