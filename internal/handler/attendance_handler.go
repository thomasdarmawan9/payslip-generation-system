package handler

import (
	"net/http"
	atDTO "payslip-generation-system/internal/dto/attendance"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

// SubmitAttendanceHandler godoc
// @Summary      Submit attendance (weekday only)
// @Description  Users can submit one attendance per day. Weekend submissions are rejected. If already submitted for the same day, response will indicate "already_exists".
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      atDTO.SubmitAttendanceRequest  true  "Submit Attendance Request"
// @Success      200      {object}  atDTO.SubmitAttendanceResponse
// @Failure      400      {object}  utils.Response[any] "Invalid request / weekend not allowed"
// @Failure      401      {object}  utils.Response[any] "Unauthorized"
// @Failure      403      {object}  utils.Response[any] "Forbidden"
// @Failure      408      {object}  utils.Response[any] "Request Process Timeout"
// @Failure      500      {object}  utils.Response[any] "Internal Server Error"
// @Router       /v1/attendance/submit [post]
func (h *Handler) SubmitAttendanceHandler(c *gin.Context) error {
	var req atDTO.SubmitAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Invalid request body",
		})
		return utils.MakeError(errorUc.BadRequest, "invalid request body")
	}

	uidAny, ok := c.Get("user_id")
	if !ok {
		h.log.Error(log.LogData{
			Description: "user_id not found in context",
		})
		return utils.MakeError(errorUc.InternalServerError, "user_id not found in context")
	}
	userID, _ := uidAny.(uint)

	row, existed, err := h.usecase.SubmitAttendance(c, userID, req.Date)
	if err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to submit attendance",
		})
		return utils.MakeError(errorUc.InternalServerError, "failed to submit attendance")
	}

	status := "created"
	if existed {
		status = "already_exists"
	}
	c.JSON(http.StatusOK, atDTO.SubmitAttendanceResponse{
		ID:     row.ID,
		UserID: row.UserID,
		Date:   row.Date.Format("2006-01-02"),
		Status: status,
	})

	return nil
}
