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

func TestSubmitAttendance_WeekendBlocked(t *testing.T) {
	u := usecase.NewForTest()

	atMock := &testm.ATRepoMock{
		CreateIfNotExistsFn: func(_ context.Context, userID uint, date time.Time) (*model.Attendance, bool, error) {
			return nil, false, nil
		},
	}
	payMock := &testm.PayRepoMock{
		HasRunOnDateFn: func(_ context.Context, date time.Time) (bool, error) { return false, nil },
	}

	// inject: AP=nil, AT=atMock, OT=nil, RB=nil, PAY=payMock, TX=Fake
	usecase.InjectForTest(u, nil, atMock, nil, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	_, _, err := u.SubmitAttendance(ctx, 1, "2025-08-17") // Minggu
	require.Error(t, err)
	require.Contains(t, err.Error(), "weekend")
}

func TestSubmitAttendance_IdempotentSameDay(t *testing.T) {
	u := usecase.NewForTest()

	atMock := &testm.ATRepoMock{
		CreateIfNotExistsFn: func(_ context.Context, userID uint, date time.Time) (*model.Attendance, bool, error) {
			return &model.Attendance{ID: 7, UserID: userID, Date: date}, true, nil
		},
	}
	payMock := &testm.PayRepoMock{
		HasRunOnDateFn: func(_ context.Context, date time.Time) (bool, error) { return false, nil },
	}

	usecase.InjectForTest(u, nil, atMock, nil, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	row, existed, err := u.SubmitAttendance(ctx, 42, "2025-08-18")
	require.NoError(t, err)
	require.True(t, existed)
	require.Equal(t, uint(42), row.UserID)
}
