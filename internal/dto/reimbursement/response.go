package reimbursement

type ReimbursementResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}
