package attendance

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	repotx "payslip-generation-system/internal/repository/tx"

	"gorm.io/gorm"
)

type Repo interface {
	CreateIfNotExists(ctx context.Context, userID uint, date time.Time) (*model.Attendance, bool, error)
}

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repo { return &repo{db: db} }

func (r *repo) CreateIfNotExists(ctx context.Context, userID uint, date time.Time) (*model.Attendance, bool, error) {
	db := repotx.GetDB(ctx, r.db)

	var existing model.Attendance
	err := db.Where("user_id = ? AND date = ?", userID, date).
		First(&existing).Error
	if err == nil {
		return &existing, true, nil // already exists
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, false, err
	}

	row := &model.Attendance{
		UserID: userID,
		Date:   date,
	}
	if err := db.Create(row).Error; err != nil {
		return nil, false, err
	}
	return row, false, nil
}
