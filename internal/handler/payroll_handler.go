// internal/handler/payroll_handler.go
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	pDTO "payslip-generation-system/internal/dto/payroll"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

// RunPayrollHandler godoc
// @Summary      Run payroll for a period (admin only)
// @Description  Processes payslips for the specified attendance period. After run, submissions in that period won't affect payslip. Can only run once per period.
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer JWT Token"
// @Param        period_id  path  int  true  "Attendance Period ID"
// @Success      200  {object}  pDTO.RunPayrollResponse
// @Failure      400  {object}  utils.Response[any] "Invalid period / already run / no working days"
// @Failure      401  {object}  utils.Response[any] "Unauthorized"
// @Failure      403  {object}  utils.Response[any] "Admin only"
// @Failure      408  {object}  utils.Response[any] "Request Process Timeout"
// @Failure      500  {object}  utils.Response[any] "Internal Server Error"
// @Router       /v1/payroll/periods/{period_id}/run [post]
func (h *Handler) RunPayrollHandler(c *gin.Context) error {
	pidStr := c.Param("period_id")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil || pid64 == 0 {
		return utils.MakeError(errorUc.BadRequest, "invalid period_id")
	}

	run, items, err := h.usecase.RunPayroll(c, uint(pid64))
	if err != nil {
		h.log.Error(log.LogData{Err: err, Description: "failed to run payroll"})
		return err
	}

	resp := pDTO.RunPayrollResponse{
		RunID:    run.ID,
		PeriodID: run.PeriodID,
		Items:    make([]pDTO.PayrollItemSummary, 0, len(items)),
	}
	for _, it := range items {
		resp.Items = append(resp.Items, pDTO.PayrollItemSummary{
			UserID:             it.UserID,
			SnapshotSalary:     fmt.Sprintf("%.2f", it.SnapshotSalary),
			WorkingDays:        it.WorkingDays,
			AttendanceDays:     it.AttendanceDays,
			OvertimeHours:      fmt.Sprintf("%.2f", it.OvertimeHours),
			BasePay:            fmt.Sprintf("%.2f", it.BasePay),
			OvertimePay:        fmt.Sprintf("%.2f", it.OvertimePay),
			ReimbursementTotal: fmt.Sprintf("%.2f", it.ReimbursementTotal),
			GrandTotal:         fmt.Sprintf("%.2f", it.GrandTotal),
		})
	}

	c.JSON(http.StatusOK, resp)
	return nil
}
