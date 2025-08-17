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

func TestCreateReimbursement_InvalidAmount(t *testing.T) {
	u := usecase.NewForTest()

	rbMock := &testm.RBRepoMock{}
	payMock := &testm.PayRepoMock{
		HasRunOnDateFn: func(_ context.Context, date time.Time) (bool, error) { return false, nil },
	}

	usecase.InjectForTest(u, nil, nil, nil, rbMock, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	_, err := u.CreateReimbursement(ctx, 7, "2025-08-18", 0, "")
	require.Error(t, err)
}

func TestCreateReimbursement_Happy(t *testing.T) {
	u := usecase.NewForTest()

	rbMock := &testm.RBRepoMock{
		CreateFn: func(_ context.Context, r *model.Reimbursement) error {
			r.ID = 11
			return nil
		},
	}
	payMock := &testm.PayRepoMock{
		HasRunOnDateFn: func(_ context.Context, date time.Time) (bool, error) { return false, nil },
	}

	usecase.InjectForTest(u, nil, nil, nil, rbMock, payMock, testm.FakeTxManager{})

	ctx := makeGinCtx()
	row, err := u.CreateReimbursement(ctx, 9, "2025-08-18", 150000, "meal")
	require.NoError(t, err)
	require.Equal(t, uint(11), row.ID)
}
