// internal/repository/reimbursement/repo.go
package reimbursement

import (
	"context"

	"payslip-generation-system/internal/model"
	repotx "payslip-generation-system/internal/repository/tx"

	"gorm.io/gorm"
)

type Repo interface {
	Create(ctx context.Context, r *model.Reimbursement) error
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db: db} }

func (r *repo) Create(ctx context.Context, m *model.Reimbursement) error {
	return repotx.GetDB(ctx, r.db).Create(m).Error
}
