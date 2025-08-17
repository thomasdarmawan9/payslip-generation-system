package usecase

import (
	"fmt"
	"math"
	"time"

	"payslip-generation-system/internal/dto/payslip"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/internal/model"
	payRepo "payslip-generation-system/internal/repository/payroll"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

func round3(v float64) float64 { return math.Round(v*100) / 100 }

func (u *usecase) GeneratePayslip(ctx *gin.Context, userID, periodID uint) (*payslip.PayslipResponse, error) {
	var pr payRepo.Repo = u.payrollRepo

	period, err := pr.GetPeriodByID(ctx, periodID)
	if err != nil {
		return nil, utils.MakeError(errorUc.BadRequest, "attendance period not found")
	}
	start := time.Date(period.StartDate.Year(), period.StartDate.Month(), period.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(period.EndDate.Year(), period.EndDate.Month(), period.EndDate.Day(), 0, 0, 0, 0, time.UTC)

	// response base
	resp := &payslip.PayslipResponse{}
	resp.Period.ID = period.ID
	resp.Period.Name = period.Name
	resp.Period.StartDate = period.StartDate.Format("2006-01-02")
	resp.Period.EndDate = period.EndDate.Format("2006-01-02")
	resp.OvertimeMultiplier = 2.0

	// Sudah run?
	run, errRun := pr.GetRunByPeriod(ctx, periodID)
	if errRun == nil {
		// gunakan snapshot payroll_items
		item, err := pr.GetPayrollItemByUser(ctx, run.ID, userID)
		if err != nil {
			// user mungkin tidak punya item (tidak ada salary/aktivitas) → tetap 0
			item = &model.PayrollItem{
				UserID:         userID,
				SnapshotSalary: 0,
				WorkingDays:    0, AttendanceDays: 0,
				WorkingHours: 0, AttendanceHours: 0,
				OvertimeHours: 0, BasePay: 0, OvertimePay: 0,
				ReimbursementTotal: 0, GrandTotal: 0,
			}
		}
		resp.SnapshotUsed = true
		resp.WorkingDays = item.WorkingDays
		resp.AttendanceDays = item.AttendanceDays
		resp.WorkingHours = item.WorkingHours
		resp.AttendanceHours = item.AttendanceHours
		// hourly dari snapshot salary / working hours (hindari div 0)
		hourly := 0.0
		if item.WorkingHours > 0 {
			hourly = item.SnapshotSalary / float64(item.WorkingHours)
		}
		resp.HourlyRate = fmt.Sprintf("%.2f", round3(hourly))
		resp.BasePay = fmt.Sprintf("%.2f", round3(item.BasePay))
		resp.OvertimeHours = fmt.Sprintf("%.2f", round3(item.OvertimeHours))
		resp.OvertimePay = fmt.Sprintf("%.2f", round3(item.OvertimePay))
		resp.SalarySnapshot = fmt.Sprintf("%.2f", round3(item.SnapshotSalary))

		// list reimburse (aman karena period terkunci)
		reims, err := pr.ListReimbursementsForUser(ctx, userID, start, end)
		if err != nil {
			return nil, utils.MakeError(errorUc.InternalServerError, "db error (reimburse list)")
		}
		sum := 0.0
		resp.Reimbursements = make([]payslip.ReimbursementLine, 0, len(reims))
		for _, r := range reims {
			sum += r.Amount
			resp.Reimbursements = append(resp.Reimbursements, payslip.ReimbursementLine{
				ID:          r.ID,
				Date:        r.Date.Format("2006-01-02"),
				Amount:      fmt.Sprintf("%.2f", round3(r.Amount)),
				Description: r.Description,
			})
		}
		resp.ReimbursementSum = fmt.Sprintf("%.2f", round3(sum))
		resp.GrandTotal = fmt.Sprintf("%.2f", round3(item.BasePay+item.OvertimePay+sum))
		return resp, nil
	}

	// Belum run → hitung on-the-fly
	workingDays := u.workingWeekdays(start, end)
	workingHours := workingDays * 8
	if workingDays <= 0 || workingHours <= 0 {
		return nil, utils.MakeError(errorUc.BadRequest, "period has no working days")
	}

	salary, err := pr.GetUserSalary(ctx, userID)
	if err != nil {
		return nil, utils.MakeError(errorUc.InternalServerError, "db error (salary)")
	}
	attDays, err := pr.GetAttendanceDaysForUser(ctx, userID, start, end)
	if err != nil {
		return nil, utils.MakeError(errorUc.InternalServerError, "db error (attendance)")
	}
	otHours, err := pr.GetOvertimeHoursForUser(ctx, userID, start, end)
	if err != nil {
		return nil, utils.MakeError(errorUc.InternalServerError, "db error (overtime)")
	}
	reims, err := pr.ListReimbursementsForUser(ctx, userID, start, end)
	if err != nil {
		return nil, utils.MakeError(errorUc.InternalServerError, "db error (reimburse list)")
	}

	attHours := attDays * 8
	hourly := 0.0
	if workingHours > 0 {
		hourly = salary / float64(workingHours)
	}
	basePay := round3(float64(attHours) * hourly)
	overtimePay := round3(otHours * (hourly * 2))
	sum := 0.0
	lines := make([]payslip.ReimbursementLine, 0, len(reims))
	for _, r := range reims {
		sum += r.Amount
		lines = append(lines, payslip.ReimbursementLine{
			ID:          r.ID,
			Date:        r.Date.Format("2006-01-02"),
			Amount:      fmt.Sprintf("%.2f", round3(r.Amount)),
			Description: r.Description,
		})
	}

	resp.SnapshotUsed = false
	resp.WorkingDays = workingDays
	resp.AttendanceDays = attDays
	resp.WorkingHours = workingHours
	resp.AttendanceHours = attHours
	resp.HourlyRate = fmt.Sprintf("%.2f", round3(hourly))
	resp.BasePay = fmt.Sprintf("%.2f", basePay)
	resp.OvertimeHours = fmt.Sprintf("%.2f", round3(otHours))
	resp.OvertimePay = fmt.Sprintf("%.2f", overtimePay)
	resp.Reimbursements = lines
	resp.ReimbursementSum = fmt.Sprintf("%.2f", round3(sum))
	resp.SalarySnapshot = fmt.Sprintf("%.2f", round3(salary))
	resp.GrandTotal = fmt.Sprintf("%.2f", round3(basePay+overtimePay+sum))
	return resp, nil
}
