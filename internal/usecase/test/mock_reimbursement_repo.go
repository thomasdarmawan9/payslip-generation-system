package test

import (
	"context"

	"payslip-generation-system/internal/model"
	rbRepo "payslip-generation-system/internal/repository/reimbursement"
)

type RBRepoMock struct {
	CreateFn func(ctx context.Context, r *model.Reimbursement) error
}

func (m *RBRepoMock) Create(ctx context.Context, r *model.Reimbursement) error {
	return m.CreateFn(ctx, r)
}

var _ rbRepo.Repo = (*RBRepoMock)(nil)
