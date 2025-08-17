package handler

import (
	"net/http"
	apDTO "payslip-generation-system/internal/dto/attendance_period"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

// CreateAttendancePeriodHandler godoc
// @Summary      Create payroll attendance period
// @Description  Admin membuat periode payroll (tidak boleh overlap, end_date >= start_date). Tanggal format YYYY-MM-DD.
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer JWT Token"
// @Param        request  body      apDTO.CreatePeriodRequest  true  "Create Payroll Period Request"
// @Success      201      {object}  apDTO.PeriodResponse
// @Failure      400      {object}  utils.Response[any] "Invalid request body / invalid dates / overlapping period"
// @Failure      401      {object}  utils.Response[any] "Unauthorized"
// @Failure      403      {object}  utils.Response[any] "Admin only"
// @Failure      408      {object}  utils.Response[any] "Request Process Timeout"
// @Failure      500      {object}  utils.Response[any] "Internal Server Error"
// @Router       /v1/payroll/periods [post]
func (h *Handler) CreateAttendancePeriodHandler(c *gin.Context) error {
	var req apDTO.CreatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Invalid request body",
		})
		return utils.MakeError(errorUc.BadRequest, "invalid request body")
	}

	row, err := h.usecase.CreateAttendancePeriod(c, req.Name, req.StartDate, req.EndDate)
	if err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to create attendance period",
		})
		return utils.MakeError(errorUc.InternalServerError, "failed to create attendance period")
	}

	h.log.Info(log.LogData{Description: "create attendance period success", Response: row})

	c.JSON(http.StatusCreated, apDTO.PeriodResponse{
		ID:        row.ID,
		Name:      row.Name,
		StartDate: row.StartDate.Format("2006-01-02"),
		EndDate:   row.EndDate.Format("2006-01-02"),
	})
	return nil
}
