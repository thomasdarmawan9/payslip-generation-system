package handler

import (
	"net/http"
	"strconv"

	psDTO "payslip-generation-system/internal/dto/payslip"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

// GeneratePayslipHandler godoc
// @Summary      Generate payslip for a period (employee)
// @Description  Generates a payslip with attendance, overtime, reimbursements and totals. If payroll already ran for the period, snapshot values are used.
// @Tags         Payslip
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer JWT Token"
// @Param        period_id  path  int  true  "Attendance Period ID"
// @Success      200  {object}  psDTO.PayslipResponse
// @Failure      400  {object}  utils.Response[any] "Invalid period / no working days"
// @Failure      401  {object}  utils.Response[any] "Unauthorized"
// @Failure      403  {object}  utils.Response[any] "Forbidden"
// @Failure      408  {object}  utils.Response[any] "Request Process Timeout"
// @Failure      500  {object}  utils.Response[any] "Internal Server Error"
// @Router       /v1/payslips/periods/{period_id} [get]
func (h *Handler) GeneratePayslipHandler(c *gin.Context) error {
	pidStr := c.Param("period_id")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil || pid64 == 0 {
		return utils.MakeError(errorUc.BadRequest, "invalid period_id")
	}

	uidAny, ok := c.Get("user_id")
	if !ok {
		return utils.MakeError(errorUc.ErrUnauthorized)
	}
	userID, _ := uidAny.(uint)

	resp, err := h.usecase.GeneratePayslip(c, userID, uint(pid64))
	if err != nil {
		h.log.Error(log.LogData{Err: err, Description: "failed to generate payslip"})
		return err
	}

	// respect timeout
	if c.Request.Context().Err() != nil {
		return nil
	}

	psDTO := psDTO.PayslipResponse{
		AttendanceDays:   resp.AttendanceDays,
		WorkingHours:     resp.WorkingHours,
		AttendanceHours:  resp.AttendanceHours,
		HourlyRate:       resp.HourlyRate,
		BasePay:          resp.BasePay,
		OvertimeHours:    resp.OvertimeHours,
		OvertimePay:      resp.OvertimePay,
		Reimbursements:   resp.Reimbursements,
		ReimbursementSum: resp.ReimbursementSum,
		SalarySnapshot:   resp.SalarySnapshot,
		GrandTotal:       resp.GrandTotal,
	}

	_ = psDTO

	c.JSON(http.StatusOK, resp)
	return nil
}
