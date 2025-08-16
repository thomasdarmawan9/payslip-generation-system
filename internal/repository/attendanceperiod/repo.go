package attendanceperiod

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	repotx "payslip-generation-system/internal/repository/tx"

	"gorm.io/gorm"
)

type Repo interface {
	Create(ctx context.Context, p *model.AttendancePeriod) error
	IsOverlapping(ctx context.Context, start, end time.Time) (bool, error)
}

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repo { return &repo{db: db} }

func (r *repo) Create(ctx context.Context, p *model.AttendancePeriod) error {
	db := repotx.GetDB(ctx, r.db)
	return db.Create(p).Error
}

func (r *repo) IsOverlapping(ctx context.Context, start, end time.Time) (bool, error) {
	db := repotx.GetDB(ctx, r.db)
	var count int64
	err := db.Model(&model.AttendancePeriod{}).
		Where("start_date <= ? AND end_date >= ?", end, start).
		Count(&count).Error
	return count > 0, err
}
