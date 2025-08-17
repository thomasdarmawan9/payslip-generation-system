package test

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	atRepo "payslip-generation-system/internal/repository/attendance"
)

type ATRepoMock struct {
	CreateIfNotExistsFn func(ctx context.Context, userID uint, date time.Time) (*model.Attendance, bool, error)
}

func (m *ATRepoMock) CreateIfNotExists(ctx context.Context, userID uint, date time.Time) (*model.Attendance, bool, error) {
	return m.CreateIfNotExistsFn(ctx, userID, date)
}

var _ atRepo.Repo = (*ATRepoMock)(nil)
