package usecase

import (
	"time"

	"github.com/gin-gonic/gin"

	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"
)

func (u *usecase) CreateAttendancePeriod(ctx *gin.Context, name, start, end string) (*model.AttendancePeriod, error) {
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, utils.MakeError(errorUc.BadRequest, "invalid start_date format (YYYY-MM-DD)")
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, utils.MakeError(errorUc.BadRequest, "invalid end_date format (YYYY-MM-DD)")
	}
	if endDate.Before(startDate) {
		return nil, utils.MakeError(errorUc.BadRequest, "end_date must be >= start_date")
	}

	txCtx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			_ = u.txManager.Rollback(txCtx)
		} else {
			if cmErr := u.txManager.Commit(txCtx); cmErr != nil {
				_ = u.txManager.Rollback(txCtx)
				err = utils.MakeError(errorUc.InternalServerError, "failed to commit transaction")
			}
		}
	}()

	overlap, err := u.apRepo.IsOverlapping(txCtx, startDate, endDate)
	if err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, utils.MakeError(errorUc.InternalServerError, "db error")
	}
	if overlap {
		return nil, utils.MakeError(errorUc.BadRequest, "period overlaps existing payroll period")
	}

	row := &model.AttendancePeriod{
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
	}
	if err := u.apRepo.Create(txCtx, row); err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to create period")
	}

	return row, nil
}

func (u *usecase) SubmitAttendance(ctx *gin.Context, userID uint, dateStr string) (*model.Attendance, bool, error) {
	// default ke "hari ini" (WIB)
	var date time.Time
	var err error
	if dateStr == "" {
		now := time.Now().In(time.FixedZone("WIB", 7*3600))
		date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		date, err = time.ParseInLocation("2006-01-02", dateStr, time.FixedZone("WIB", 7*3600))
		if err != nil {
			return nil, false, utils.MakeError(errorUc.BadRequest, "invalid date format (YYYY-MM-DD)")
		}
	}

	// Rule: tidak boleh submit weekend (Sabtu = 6, Minggu = 0; tergantung locale).
	// Go: Sunday=0 ... Saturday=6
	wd := date.Weekday()
	if wd == time.Saturday || wd == time.Sunday {
		return nil, false, utils.MakeError(errorUc.BadRequest, "cannot submit attendance on weekend")
	}

	txCtx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, false, utils.MakeError(errorUc.InternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			_ = u.txManager.Rollback(txCtx)
		} else {
			if cmErr := u.txManager.Commit(txCtx); cmErr != nil {
				_ = u.txManager.Rollback(txCtx)
				err = utils.MakeError(errorUc.InternalServerError, "failed to commit transaction")
			}
		}
	}()

	row, existed, err := u.atRepo.CreateIfNotExists(txCtx, userID, date)
	if err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, false, utils.MakeError(errorUc.InternalServerError, "db error")
	}

	return row, existed, nil
}
