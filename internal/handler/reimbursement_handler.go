// internal/handler/reimbursement_handler.go
package handler

import (
	"fmt"
	"net/http"

	rbDTO "payslip-generation-system/internal/dto/reimbursement"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

// CreateReimbursementHandler godoc
// @Summary      Create reimbursement
// @Description  Create a reimbursement with amount and optional description.
// @Tags         Reimbursement
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      rbDTO.CreateReimbursementRequest  true  "Create Reimbursement Request"
// @Success      201      {object}  rbDTO.ReimbursementResponse
// @Failure      400      {object}  utils.Response[any] "Invalid request / amount <= 0"
// @Failure      401      {object}  utils.Response[any] "Unauthorized"
// @Failure      403      {object}  utils.Response[any] "Forbidden"
// @Failure      408      {object}  utils.Response[any] "Request Process Timeout"
// @Failure      500      {object}  utils.Response[any] "Internal Server Error"
// @Router       /v1/reimbursements [post]
func (h *Handler) CreateReimbursementHandler(c *gin.Context) error {
	var req rbDTO.CreateReimbursementRequest
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

	row, err := h.usecase.CreateReimbursement(c, userID, req.Date, req.Amount, req.Description)
	if err != nil {
		h.log.Error(log.LogData{Err: err, Description: "Failed to create reimbursement"})
		return err
	}

	c.JSON(http.StatusCreated, rbDTO.ReimbursementResponse{
		ID:          row.ID,
		UserID:      row.UserID,
		Date:        row.Date.Format("2006-01-02"),
		Amount:      fmt.Sprintf("%.2f", row.Amount),
		Description: row.Description,
	})
	return nil
}
