// internal/usecase/overtime_reimbursement_usecase.go
package usecase

import (
	"time"

	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

func (u *usecase) SubmitOvertime(ctx *gin.Context, userID uint, dateStr string, hours float64) (*model.Overtime, bool, error) {
	// Validasi jam
	if hours <= 0 || hours > 3 {
		return nil, false, utils.MakeError(errorUc.BadRequest, "hours must be > 0 and <= 3")
	}

	// Parse tanggal (default today WIB)
	loc := time.FixedZone("WIB", 7*3600)
	var date time.Time
	var err error
	if dateStr == "" {
		now := time.Now().In(loc)
		date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	} else {
		date, err = time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			return nil, false, utils.MakeError(errorUc.BadRequest, "invalid date format (YYYY-MM-DD)")
		}
	}

	// Harus diajukan setelah jam kerja selesai (>= 17:00 WIB) kalau tanggal = hari ini
	now := time.Now().In(loc)
	if now.Year() == date.Year() && now.YearDay() == date.YearDay() {
		after5pm := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, loc)
		if now.Before(after5pm) {
			return nil, false, utils.MakeError(errorUc.BadRequest, "overtime can only be submitted after 17:00 WIB")
		}
	}
	// Catatan: Overtime bisa diambil hari apa pun (weekend allowed) → tidak ada cek weekend.

	txCtx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, false, utils.MakeError(errorUc.InternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			_ = u.txManager.Rollback(txCtx)
		} else if cmErr := u.txManager.Commit(txCtx); cmErr != nil {
			_ = u.txManager.Rollback(txCtx)
			err = utils.MakeError(errorUc.InternalServerError, "failed to commit transaction")
		}
	}()

	row, existed, err := u.otRepo.CreateIfNotExists(txCtx, userID, date, hours)
	if err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, false, utils.MakeError(errorUc.InternalServerError, "db error")
	}
	return row, existed, nil
}

func (u *usecase) CreateReimbursement(ctx *gin.Context, userID uint, dateStr string, amount float64, description string) (*model.Reimbursement, error) {
	if amount <= 0 {
		return nil, utils.MakeError(errorUc.BadRequest, "amount must be > 0")
	}
	loc := time.FixedZone("WIB", 7*3600)
	var date time.Time
	var err error
	if dateStr == "" {
		now := time.Now().In(loc)
		date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	} else {
		date, err = time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			return nil, utils.MakeError(errorUc.BadRequest, "invalid date format (YYYY-MM-DD)")
		}
	}

	txCtx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			_ = u.txManager.Rollback(txCtx)
		} else if cmErr := u.txManager.Commit(txCtx); cmErr != nil {
			_ = u.txManager.Rollback(txCtx)
			err = utils.MakeError(errorUc.InternalServerError, "failed to commit transaction")
		}
	}()

	row := &model.Reimbursement{
		UserID:      userID,
		Date:        date,
		Amount:      amount,
		Description: description,
	}
	if err := u.rbRepo.Create(txCtx, row); err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, utils.MakeError(errorUc.InternalServerError, "db error")
	}
	return row, nil
}
