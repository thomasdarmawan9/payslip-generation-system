// internal/usecase/payroll_usecase.go
package usecase

import (
	"math"
	"time"

	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/internal/model"
	payRepo "payslip-generation-system/internal/repository/payroll"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

func (u *usecase) workingWeekdays(start, end time.Time) int {
	d := 0
	for cur := start; !cur.After(end); cur = cur.AddDate(0, 0, 1) {
		switch cur.Weekday() {
		case time.Saturday, time.Sunday:
			continue
		default:
			d++
		}
	}
	return d
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func (u *usecase) RunPayroll(ctx *gin.Context, periodID uint) (*model.PayrollRun, []*model.PayrollItem, error) {
	// Repos
	var pr payRepo.Repo
	if u.payrollRepo == nil {
		// injected via provider
	}
	pr = u.payrollRepo

	period, err := pr.GetPeriodByID(ctx, periodID)
	if err != nil {
		return nil, nil, utils.MakeError(errorUc.BadRequest, "attendance period not found")
	}

	// only once per period
	exists, err := pr.HasRunForPeriod(ctx, periodID)
	if err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "db error")
	}
	if exists {
		return nil, nil, utils.MakeError(errorUc.BadRequest, "payroll has already been run for this period")
	}

	start := time.Date(period.StartDate.Year(), period.StartDate.Month(), period.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(period.EndDate.Year(), period.EndDate.Month(), period.EndDate.Day(), 0, 0, 0, 0, time.UTC)

	workingDays := u.workingWeekdays(start, end)
	workingHours := workingDays * 8
	if workingDays <= 0 || workingHours <= 0 {
		return nil, nil, utils.MakeError(errorUc.BadRequest, "period has no working days")
	}

	// aggregates
	attDays, err := pr.GetAttendanceDaysByUser(ctx, start, end)
	if err != nil {
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "db error (attendance agg)")
	}
	otHours, err := pr.GetOvertimeHoursByUser(ctx, start, end)
	if err != nil {
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "db error (overtime agg)")
	}
	rbTotals, err := pr.GetReimbTotalByUser(ctx, start, end)
	if err != nil {
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "db error (reimburse agg)")
	}
	salaries, err := pr.GetUserSalaries(ctx)
	if err != nil {
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "db error (salaries)")
	}

	// build items untuk semua user yang punya attendance/overtime/reimburse ataupun punya salary
	userSet := map[uint]struct{}{}
	for uid := range salaries {
		userSet[uid] = struct{}{}
	}
	for uid := range attDays {
		userSet[uid] = struct{}{}
	}
	for uid := range otHours {
		userSet[uid] = struct{}{}
	}
	for uid := range rbTotals {
		userSet[uid] = struct{}{}
	}

	items := make([]*model.PayrollItem, 0, len(userSet))
	for uid := range userSet {
		sal := salaries[uid] 
		att := attDays[uid]
		ot := otHours[uid]
		rbt := rbTotals[uid]

		attHours := att * 8
		hourly := 0.0
		if workingHours > 0 {
			hourly = sal / float64(workingHours)
		}
		basePay := round2(float64(attHours) * hourly)
		overtimePay := round2(ot * (hourly * 2))
		total := round2(basePay + overtimePay + rbt)

		items = append(items, &model.PayrollItem{
			UserID:             uid,
			SnapshotSalary:     round2(sal),
			WorkingDays:        workingDays,
			AttendanceDays:     att,
			WorkingHours:       workingHours,
			AttendanceHours:    attHours,
			OvertimeHours:      round2(ot),
			BasePay:            basePay,
			OvertimePay:        overtimePay,
			ReimbursementTotal: round2(rbt),
			GrandTotal:         total,
		})
	}

	txCtx, err := u.txManager.Begin(ctx)
	if err != nil {
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			_ = u.txManager.Rollback(txCtx)
		} else if cmErr := u.txManager.Commit(txCtx); cmErr != nil {
			_ = u.txManager.Rollback(txCtx)
			err = utils.MakeError(errorUc.InternalServerError, "failed to commit transaction")
		}
	}()

	run := &model.PayrollRun{
		PeriodID: periodID,
		RunAt:    time.Now().UTC(),
	}
	if err := pr.CreateRun(txCtx, run, items); err != nil {
		u.log.Error(log.LogData{Err: err})
		return nil, nil, utils.MakeError(errorUc.InternalServerError, "failed to persist payroll")
	}

	return run, items, nil
}
