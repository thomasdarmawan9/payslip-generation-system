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

func TestCreateAttendancePeriod_HappyPath(t *testing.T) {
	u := usecase.NewForTest()

	apMock := &testm.APRepoMock{
		OverlapFn: func(_ context.Context, start, end time.Time) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, p *model.AttendancePeriod) error {
			p.ID = 123
			return nil
		},
	}

	usecase.InjectForTest(
		u,
		apMock,                // attendance period repo
		nil,                   // attendance repo
		nil,                   // overtime repo
		nil,                   // reimbursement repo
		nil,                   // payroll repo
		testm.FakeTxManager{}, // tx manager
	)

	ctx := makeGinCtx()
	row, err := u.CreateAttendancePeriod(ctx, "Agustus 2025", "2025-08-01", "2025-08-31")
	require.NoError(t, err)
	require.Equal(t, uint(123), row.ID)
	require.Equal(t, "Agustus 2025", row.Name)
}

func TestCreateAttendancePeriod_Overlap(t *testing.T) {
	u := usecase.NewForTest()

	apMock := &testm.APRepoMock{
		OverlapFn: func(_ context.Context, start, end time.Time) (bool, error) {
			return true, nil
		},
		CreateFn: func(_ context.Context, p *model.AttendancePeriod) error {
			return nil
		},
	}

	usecase.InjectForTest(u, apMock, nil, nil, nil, nil, testm.FakeTxManager{})

	ctx := makeGinCtx()
	row, err := u.CreateAttendancePeriod(ctx, "X", "2025-08-01", "2025-08-31")
	require.Error(t, err)
	require.Nil(t, row)
	require.Contains(t, err.Error(), "overlap")
}
