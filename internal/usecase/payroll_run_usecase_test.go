package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"payslip-generation-system/internal/model"
	"payslip-generation-system/internal/usecase"
	testm "payslip-generation-system/internal/usecase/test"
)

func TestRunPayroll_OnlyOnce(t *testing.T) {
	u := usecase.NewForTest()
	payMock := &testm.PayRepoMock{
		GetPeriodByIDFn: func(_ context.Context, id uint) (*model.AttendancePeriod, error) {
			return &model.AttendancePeriod{
				ID:        id,
				StartDate: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		HasRunForPeriodFn: func(_ context.Context, periodID uint) (bool, error) { return true, nil },
	}
	usecase.InjectForTest(u, nil, nil, nil, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	run, items, err := u.RunPayroll(ctx, 1)
	require.Error(t, err)
	require.Nil(t, run)
	require.Nil(t, items)
}

func TestRunPayroll_CalcNumbers(t *testing.T) {
	u := usecase.NewForTest()
	payMock := &testm.PayRepoMock{
		GetPeriodByIDFn: func(_ context.Context, id uint) (*model.AttendancePeriod, error) {
			return &model.AttendancePeriod{
				ID:        id,
				StartDate: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		HasRunForPeriodFn:         func(_ context.Context, periodID uint) (bool, error) { return false, nil },
		GetAttendanceDaysByUserFn: func(_ context.Context, s, e time.Time) (map[uint]int, error) { return map[uint]int{7: 20}, nil },
		GetOvertimeHoursByUserFn:  func(_ context.Context, s, e time.Time) (map[uint]float64, error) { return map[uint]float64{7: 5}, nil },
		GetReimbTotalByUserFn: func(_ context.Context, s, e time.Time) (map[uint]float64, error) {
			return map[uint]float64{7: 100000}, nil
		},
		GetUserSalariesFn: func(_ context.Context) (map[uint]float64, error) { return map[uint]float64{7: 7000000}, nil },
		CreateRunFn: func(_ context.Context, run *model.PayrollRun, items []*model.PayrollItem) error {
			run.ID = 99
			return nil
		},
	}
	usecase.InjectForTest(u, nil, nil, nil, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	run, items, err := u.RunPayroll(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, uint(99), run.ID)
	require.Len(t, items, 1)
	require.Equal(t, uint(7), items[0].UserID)
	require.Greater(t, items[0].GrandTotal, 0.0)
}
