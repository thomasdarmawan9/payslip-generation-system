package payslip

type ReimbursementLine struct {
	ID          uint   `json:"id"`
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

type PayslipResponse struct {
	Period struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"period"`

	SnapshotUsed bool `json:"snapshot_used"` // true jika payroll sudah run

	// Breakdown attendance / base pay
	WorkingDays     int    `json:"working_days"`
	AttendanceDays  int    `json:"attendance_days"`
	WorkingHours    int    `json:"working_hours"`
	AttendanceHours int    `json:"attendance_hours"`
	HourlyRate      string `json:"hourly_rate"`
	BasePay         string `json:"base_pay"`

	// Overtime breakdown
	OvertimeHours      string  `json:"overtime_hours"`
	OvertimeMultiplier float64 `json:"overtime_multiplier"` // 2.0
	OvertimePay        string  `json:"overtime_pay"`

	// Reimbursements
	Reimbursements   []ReimbursementLine `json:"reimbursements"`
	ReimbursementSum string              `json:"reimbursement_sum"`

	// Totals
	SalarySnapshot string `json:"salary_snapshot"`
	GrandTotal     string `json:"grand_total"`
}
