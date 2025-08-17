package payroll

import (
	"context"
	"time"

	"payslip-generation-system/internal/model"
	repotx "payslip-generation-system/internal/repository/tx"

	"gorm.io/gorm"
)

type Repo interface {
	HasRunForPeriod(ctx context.Context, periodID uint) (bool, error)
	CreateRun(ctx context.Context, run *model.PayrollRun, items []*model.PayrollItem) error

	// Aggregations
	GetAttendanceDaysByUser(ctx context.Context, start, end time.Time) (map[uint]int, error)
	GetOvertimeHoursByUser(ctx context.Context, start, end time.Time) (map[uint]float64, error)
	GetReimbTotalByUser(ctx context.Context, start, end time.Time) (map[uint]float64, error)

	GetUserSalaries(ctx context.Context) (map[uint]float64, error)

	// Period lookup
	GetPeriodByID(ctx context.Context, id uint) (*model.AttendancePeriod, error)

	// Check if a date falls into a period that already has payroll run (for locking)
	HasRunOnDate(ctx context.Context, date time.Time) (bool, error)

	// Payslip related methods
	GetPayrollItemByUser(ctx context.Context, runID uint, userID uint) (*model.PayrollItem, error)
	GetRunByPeriod(ctx context.Context, periodID uint) (*model.PayrollRun, error)
	GetUserSalary(ctx context.Context, userID uint) (float64, error)
	GetAttendanceDaysForUser(ctx context.Context, userID uint, start, end time.Time) (int, error)
	GetOvertimeHoursForUser(ctx context.Context, userID uint, start, end time.Time) (float64, error)
	ListReimbursementsForUser(ctx context.Context, userID uint, start, end time.Time) ([]model.Reimbursement, error)
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repo { return &repo{db: db} }

func (r *repo) HasRunForPeriod(ctx context.Context, periodID uint) (bool, error) {
	db := repotx.GetDB(ctx, r.db)
	var count int64
	if err := db.Model(&model.PayrollRun{}).Where("period_id = ?", periodID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repo) CreateRun(ctx context.Context, run *model.PayrollRun, items []*model.PayrollItem) error {
	db := repotx.GetDB(ctx, r.db)
	if err := db.Create(run).Error; err != nil {
		return err
	}
	for _, it := range items {
		it.PayrollRunID = run.ID
	}
	return db.Create(&items).Error
}

func (r *repo) GetAttendanceDaysByUser(ctx context.Context, start, end time.Time) (map[uint]int, error) {
	db := repotx.GetDB(ctx, r.db)
	type row struct {
		UserID uint
		Count  int
	}
	var rows []row
	if err := db.
		Table((model.Attendance{}).TableName()).
		Select("user_id, COUNT(*) as count").
		Where("date BETWEEN ? AND ?", start, end).
		Group("user_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint]int, len(rows))
	for _, r := range rows {
		out[r.UserID] = r.Count
	}
	return out, nil
}

func (r *repo) GetOvertimeHoursByUser(ctx context.Context, start, end time.Time) (map[uint]float64, error) {
	db := repotx.GetDB(ctx, r.db)
	type row struct {
		UserID uint
		Hours  float64
	}
	var rows []row
	if err := db.
		Table((model.Overtime{}).TableName()).
		Select("user_id, COALESCE(SUM(hours),0) as hours").
		Where("date BETWEEN ? AND ?", start, end).
		Group("user_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint]float64, len(rows))
	for _, r := range rows {
		out[r.UserID] = r.Hours
	}
	return out, nil
}

func (r *repo) GetReimbTotalByUser(ctx context.Context, start, end time.Time) (map[uint]float64, error) {
	db := repotx.GetDB(ctx, r.db)
	type row struct {
		UserID uint
		Total  float64
	}
	var rows []row
	if err := db.
		Table((model.Reimbursement{}).TableName()).
		Select("user_id, COALESCE(SUM(amount),0) as total").
		Where("date BETWEEN ? AND ?", start, end).
		Group("user_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint]float64, len(rows))
	for _, r := range rows {
		out[r.UserID] = r.Total
	}
	return out, nil
}

func (r *repo) GetUserSalaries(ctx context.Context) (map[uint]float64, error) {
	db := repotx.GetDB(ctx, r.db)
	type row struct {
		ID     uint
		Salary float64
	}
	var rows []row
	if err := db.
		Table((model.User{}).TableName()).
		Select("id, salary").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint]float64, len(rows))
	for _, r := range rows {
		out[r.ID] = r.Salary
	}
	return out, nil
}

func (r *repo) GetPeriodByID(ctx context.Context, id uint) (*model.AttendancePeriod, error) {
	db := repotx.GetDB(ctx, r.db)
	var p model.AttendancePeriod
	if err := db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) HasRunOnDate(ctx context.Context, date time.Time) (bool, error) {
	db := repotx.GetDB(ctx, r.db)
	// payroll_runs join attendance_periods; cek apakah date berada dalam period yang sudah di-run
	type row struct{ Count int64 }
	var c int64
	err := db.Table((model.PayrollRun{}).TableName()+" pr").
		Joins("JOIN "+(model.AttendancePeriod{}).TableName()+" ap ON ap.id = pr.period_id").
		Where("? BETWEEN ap.start_date AND ap.end_date", date).
		Count(&c).Error
	return c > 0, err
}

// Tambahan method untuk payslip
func (r *repo) GetPayrollItemByUser(ctx context.Context, runID uint, userID uint) (*model.PayrollItem, error) {
	db := repotx.GetDB(ctx, r.db)
	var it model.PayrollItem
	if err := db.Where("payroll_run_id = ? AND user_id = ?", runID, userID).First(&it).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *repo) GetRunByPeriod(ctx context.Context, periodID uint) (*model.PayrollRun, error) {
	db := repotx.GetDB(ctx, r.db)
	var run model.PayrollRun
	if err := db.Where("period_id = ?", periodID).First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *repo) GetUserSalary(ctx context.Context, userID uint) (float64, error) {
	db := repotx.GetDB(ctx, r.db)
	type row struct {
		Salary float64
	}
	var rw row
	if err := db.
		Table((model.User{}).TableName()).
		Select("salary").
		Where("id = ?", userID).
		Scan(&rw).Error; err != nil {
		return 0, err
	}
	return rw.Salary, nil
}

func (r *repo) GetAttendanceDaysForUser(ctx context.Context, userID uint, start, end time.Time) (int, error) {
	db := repotx.GetDB(ctx, r.db)
	var c int64
	if err := db.
		Table((model.Attendance{}).TableName()).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *repo) GetOvertimeHoursForUser(ctx context.Context, userID uint, start, end time.Time) (float64, error) {
	db := repotx.GetDB(ctx, r.db)
	type row struct{ Hours float64 }
	var rw row
	if err := db.
		Table((model.Overtime{}).TableName()).
		Select("COALESCE(SUM(hours),0) as hours").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Scan(&rw).Error; err != nil {
		return 0, err
	}
	return rw.Hours, nil
}

func (r *repo) ListReimbursementsForUser(ctx context.Context, userID uint, start, end time.Time) ([]model.Reimbursement, error) {
	db := repotx.GetDB(ctx, r.db)
	var rows []model.Reimbursement
	if err := db.
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Order("date ASC, id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
