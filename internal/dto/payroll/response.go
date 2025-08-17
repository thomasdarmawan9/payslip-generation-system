// internal/dto/payroll/response.go
package payroll

type RunPayrollResponse struct {
	RunID    uint                 `json:"run_id"`
	PeriodID uint                 `json:"period_id"`
	Items    []PayrollItemSummary `json:"items"`
}

type PayrollItemSummary struct {
	UserID             uint   `json:"user_id"`
	SnapshotSalary     string `json:"snapshot_salary"`
	WorkingDays        int    `json:"working_days"`
	AttendanceDays     int    `json:"attendance_days"`
	OvertimeHours      string `json:"overtime_hours"`
	BasePay            string `json:"base_pay"`
	OvertimePay        string `json:"overtime_pay"`
	ReimbursementTotal string `json:"reimbursement_total"`
	GrandTotal         string `json:"grand_total"`
}
