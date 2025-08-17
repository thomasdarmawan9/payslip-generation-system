package test

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	apRepo "payslip-generation-system/internal/repository/attendanceperiod"
)

type APRepoMock struct {
	OverlapFn func(ctx context.Context, start, end time.Time) (bool, error)
	CreateFn  func(ctx context.Context, p *model.AttendancePeriod) error
}

func (m *APRepoMock) IsOverlapping(ctx context.Context, start, end time.Time) (bool, error) {
	return m.OverlapFn(ctx, start, end)
}
func (m *APRepoMock) Create(ctx context.Context, p *model.AttendancePeriod) error {
	return m.CreateFn(ctx, p)
}

var _ apRepo.Repo = (*APRepoMock)(nil)
