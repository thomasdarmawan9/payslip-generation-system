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

func TestSubmitOvertime_TooManyHours(t *testing.T) {
	u := usecase.NewForTest()

	otMock := &testm.OTRepoMock{} // takkan dipanggil karena validasi jam
	payMock := &testm.PayRepoMock{
		HasRunOnDateFn: func(_ context.Context, date time.Time) (bool, error) { return false, nil },
	}

	usecase.InjectForTest(u, nil, nil, otMock, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	_, _, err := u.SubmitOvertime(ctx, 1, "2025-08-18", 3.5)
	require.Error(t, err)
}

func TestSubmitOvertime_PastDate_OK(t *testing.T) {
	u := usecase.NewForTest()

	otMock := &testm.OTRepoMock{
		CreateIfNotExistsFn: func(_ context.Context, userID uint, date time.Time, hours float64) (*model.Overtime, bool, error) {
			return &model.Overtime{ID: 1, UserID: userID, Date: date, Hours: hours}, false, nil
		},
	}
	payMock := &testm.PayRepoMock{
		HasRunOnDateFn: func(_ context.Context, date time.Time) (bool, error) { return false, nil },
	}

	usecase.InjectForTest(u, nil, nil, otMock, nil, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	row, existed, err := u.SubmitOvertime(ctx, 2, "2025-08-18", 2)
	require.NoError(t, err)
	require.False(t, existed)
	require.Equal(t, uint(2), row.UserID)
}
