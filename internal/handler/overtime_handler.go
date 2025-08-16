// internal/handler/overtime_handler.go
package handler

import (
	"fmt"
	"net/http"

	otDTO "payslip-generation-system/internal/dto/overtime"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

// SubmitOvertimeHandler godoc
// @Summary      Submit overtime
// @Description  Submit overtime hours (<= 3h). Only allowed after 17:00 WIB if submitting for today. Weekend allowed.
// @Tags         Overtime
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      otDTO.SubmitOvertimeRequest  true  "Submit Overtime Request"
// @Success      200      {object}  otDTO.SubmitOvertimeResponse
// @Failure      400      {object}  utils.Response[any] "Invalid request / hours > 3 / before 17:00"
// @Failure      401      {object}  utils.Response[any] "Unauthorized"
// @Failure      403      {object}  utils.Response[any] "Forbidden"
// @Failure      408      {object}  utils.Response[any] "Request Process Timeout"
// @Failure      500      {object}  utils.Response[any] "Internal Server Error"
// @Router       /v1/overtime/submit [post]
func (h *Handler) SubmitOvertimeHandler(c *gin.Context) error {
	var req otDTO.SubmitOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(log.LogData{Err: err, Description: "Invalid request body"})
		return utils.MakeError(errorUc.BadRequest, "invalid request body")
	}

	uidAny, ok := c.Get("user_id")
	if !ok {
		h.log.Error(log.LogData{Description: "user_id not found in context"})
		return utils.MakeError(errorUc.ErrUnauthorized)
	}
	userID, _ := uidAny.(uint)

	row, existed, err := h.usecase.SubmitOvertime(c, userID, req.Date, req.Hours)
	if err != nil {
		h.log.Error(log.LogData{Err: err, Description: "Failed to submit overtime"})
		return err
	}

	status := "created"
	if existed {
		status = "already_exists"
	}
	c.JSON(http.StatusOK, otDTO.SubmitOvertimeResponse{
		ID:     row.ID,
		UserID: row.UserID,
		Date:   row.Date.Format("2006-01-02"),
		Hours:  fmt.Sprintf("%.2f", row.Hours),
		Status: status,
	})
	return nil
}
