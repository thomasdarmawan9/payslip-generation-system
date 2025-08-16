package overtime

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	repotx "payslip-generation-system/internal/repository/tx"

	"gorm.io/gorm"
)

type Repo interface {
	CreateIfNotExists(ctx context.Context, userID uint, date time.Time, hours float64) (*model.Overtime, bool, error)
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db: db} }

func (r *repo) CreateIfNotExists(ctx context.Context, userID uint, date time.Time, hours float64) (*model.Overtime, bool, error) {
	db := repotx.GetDB(ctx, r.db)

	var existing model.Overtime
	if err := db.Where("user_id = ? AND date = ?", userID, date).First(&existing).Error; err == nil {
		return &existing, true, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, false, err
	}

	row := &model.Overtime{UserID: userID, Date: date, Hours: hours}
	if err := db.Create(row).Error; err != nil {
		return nil, false, err
	}
	return row, false, nil
}
