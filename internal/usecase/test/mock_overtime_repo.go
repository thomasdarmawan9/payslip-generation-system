package test

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	otRepo "payslip-generation-system/internal/repository/overtime"
)

type OTRepoMock struct {
	CreateIfNotExistsFn func(ctx context.Context, userID uint, date time.Time, hours float64) (*model.Overtime, bool, error)
}

func (m *OTRepoMock) CreateIfNotExists(ctx context.Context, userID uint, date time.Time, hours float64) (*model.Overtime, bool, error) {
	return m.CreateIfNotExistsFn(ctx, userID, date, hours)
}

var _ otRepo.Repo = (*OTRepoMock)(nil)
