package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"payslip-generation-system/internal/model"
	"payslip-generation-system/internal/usecase"
	testm "payslip-generation-system/internal/usecase/test"
)

func TestGeneratePayslip_UsesSnapshot(t *testing.T) {
	u := usecase.NewForTest()
	payMock := &testm.PayRepoMock{
		GetPeriodByIDFn: func(_ context.Context, id uint) (*model.AttendancePeriod, error) {
			return &model.AttendancePeriod{
				ID: id, Name: "Aug 2025",
				StartDate: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		GetRunByPeriodFn: func(_ context.Context, periodID uint) (*model.PayrollRun, error) {
			return &model.PayrollRun{ID: 5, PeriodID: periodID}, nil
		},
		GetPayrollItemByUserFn: func(_ context.Context, runID uint, userID uint) (*model.PayrollItem, error) {
			return &model.PayrollItem{
				PayrollRunID: runID, UserID: userID,
				SnapshotSalary: 7000000, WorkingDays: 23, WorkingHours: 184,
				AttendanceDays: 20, AttendanceHours: 160,
				OvertimeHours: 5, BasePay: 6000000, OvertimePay: 380000,
				ReimbursementTotal: 100000, GrandTotal: 6480000,
			}, nil
		},
		ListReimbursementsForUserFn: func(_ context.Context, uid uint, s, e time.Time) ([]model.Reimbursement, error) {
			return []model.Reimbursement{{ID: 1, UserID: uid, Date: s.AddDate(0, 0, 2), Amount: 100000, Description: "meal"}}, nil
		},
	}
	usecase.InjectForTest(u, nil, nil, nil, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	resp, err := u.GeneratePayslip(ctx, 7, 1)
	require.NoError(t, err)
	require.True(t, resp.SnapshotUsed)
	require.Equal(t, "6480000.00", resp.GrandTotal)
}

func TestGeneratePayslip_LiveCalc(t *testing.T) {
	u := usecase.NewForTest()
	payMock := &testm.PayRepoMock{
		GetPeriodByIDFn: func(_ context.Context, id uint) (*model.AttendancePeriod, error) {
			return &model.AttendancePeriod{
				ID:        id,
				StartDate: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		GetRunByPeriodFn:           func(_ context.Context, pid uint) (*model.PayrollRun, error) { return nil, gorm.ErrRecordNotFound },
		GetUserSalaryFn:            func(_ context.Context, uid uint) (float64, error) { return 7000000, nil },
		GetAttendanceDaysForUserFn: func(_ context.Context, uid uint, s, e time.Time) (int, error) { return 20, nil },
		GetOvertimeHoursForUserFn:  func(_ context.Context, uid uint, s, e time.Time) (float64, error) { return 5, nil },
		ListReimbursementsForUserFn: func(_ context.Context, uid uint, s, e time.Time) ([]model.Reimbursement, error) {
			return []model.Reimbursement{{ID: 1, UserID: uid, Date: s.AddDate(0, 0, 2), Amount: 12345}}, nil
		},
	}
	usecase.InjectForTest(u, nil, nil, nil, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	resp, err := u.GeneratePayslip(ctx, 7, 1)
	require.NoError(t, err)
	require.False(t, resp.SnapshotUsed)
	require.NotEmpty(t, resp.GrandTotal)
}
