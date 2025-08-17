package test

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	payRepo "payslip-generation-system/internal/repository/payroll"
)

type PayRepoMock struct {
	// run & period
	HasRunForPeriodFn func(ctx context.Context, periodID uint) (bool, error)
	CreateRunFn       func(ctx context.Context, run *model.PayrollRun, items []*model.PayrollItem) error
	GetPeriodByIDFn   func(ctx context.Context, id uint) (*model.AttendancePeriod, error)
	GetRunByPeriodFn  func(ctx context.Context, periodID uint) (*model.PayrollRun, error)

	// aggs & salary
	GetAttendanceDaysByUserFn func(ctx context.Context, start, end time.Time) (map[uint]int, error)
	GetOvertimeHoursByUserFn  func(ctx context.Context, start, end time.Time) (map[uint]float64, error)
	GetReimbTotalByUserFn     func(ctx context.Context, start, end time.Time) (map[uint]float64, error)
	GetUserSalariesFn         func(ctx context.Context) (map[uint]float64, error)
	GetUserSalaryFn           func(ctx context.Context, userID uint) (float64, error)

	// per-user
	GetPayrollItemByUserFn      func(ctx context.Context, runID uint, userID uint) (*model.PayrollItem, error)
	GetAttendanceDaysForUserFn  func(ctx context.Context, userID uint, start, end time.Time) (int, error)
	GetOvertimeHoursForUserFn   func(ctx context.Context, userID uint, start, end time.Time) (float64, error)
	ListReimbursementsForUserFn func(ctx context.Context, userID uint, start, end time.Time) ([]model.Reimbursement, error)

	// lock
	HasRunOnDateFn func(ctx context.Context, date time.Time) (bool, error)
}

func (m *PayRepoMock) HasRunForPeriod(ctx context.Context, periodID uint) (bool, error) {
	return m.HasRunForPeriodFn(ctx, periodID)
}
func (m *PayRepoMock) CreateRun(ctx context.Context, run *model.PayrollRun, items []*model.PayrollItem) error {
	return m.CreateRunFn(ctx, run, items)
}
func (m *PayRepoMock) GetWorkingWeekdays(ctx context.Context, start, end time.Time) (int, error) {
	// tidak digunakan; usecase hitung sendiri
	return 0, nil
}
func (m *PayRepoMock) GetAttendanceDaysByUser(ctx context.Context, start, end time.Time) (map[uint]int, error) {
	return m.GetAttendanceDaysByUserFn(ctx, start, end)
}
func (m *PayRepoMock) GetOvertimeHoursByUser(ctx context.Context, start, end time.Time) (map[uint]float64, error) {
	return m.GetOvertimeHoursByUserFn(ctx, start, end)
}
func (m *PayRepoMock) GetReimbTotalByUser(ctx context.Context, start, end time.Time) (map[uint]float64, error) {
	return m.GetReimbTotalByUserFn(ctx, start, end)
}
func (m *PayRepoMock) GetUserSalaries(ctx context.Context) (map[uint]float64, error) {
	return m.GetUserSalariesFn(ctx)
}
func (m *PayRepoMock) GetPeriodByID(ctx context.Context, id uint) (*model.AttendancePeriod, error) {
	return m.GetPeriodByIDFn(ctx, id)
}
func (m *PayRepoMock) HasRunOnDate(ctx context.Context, date time.Time) (bool, error) {
	return m.HasRunOnDateFn(ctx, date)
}
func (m *PayRepoMock) GetRunByPeriod(ctx context.Context, periodID uint) (*model.PayrollRun, error) {
	return m.GetRunByPeriodFn(ctx, periodID)
}
func (m *PayRepoMock) GetPayrollItemByUser(ctx context.Context, runID uint, userID uint) (*model.PayrollItem, error) {
	return m.GetPayrollItemByUserFn(ctx, runID, userID)
}
func (m *PayRepoMock) GetUserSalary(ctx context.Context, userID uint) (float64, error) {
	return m.GetUserSalaryFn(ctx, userID)
}
func (m *PayRepoMock) GetAttendanceDaysForUser(ctx context.Context, userID uint, start, end time.Time) (int, error) {
	return m.GetAttendanceDaysForUserFn(ctx, userID, start, end)
}
func (m *PayRepoMock) GetOvertimeHoursForUser(ctx context.Context, userID uint, start, end time.Time) (float64, error) {
	return m.GetOvertimeHoursForUserFn(ctx, userID, start, end)
}
func (m *PayRepoMock) ListReimbursementsForUser(ctx context.Context, userID uint, start, end time.Time) ([]model.Reimbursement, error) {
	return m.ListReimbursementsForUserFn(ctx, userID, start, end)
}

var _ payRepo.Repo = (*PayRepoMock)(nil)
