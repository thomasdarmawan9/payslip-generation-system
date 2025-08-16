package handler

import (
	"net/http"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"
	"strings"

	"github.com/gin-gonic/gin"

	authDTO "payslip-generation-system/internal/dto/auth"
	errorUc "payslip-generation-system/internal/error"
)

// RegisterUserHandler godoc
// @Summary      Register User
// @Description  Register new user with required fields
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      authDTO.RegisterUserRequest  true  "Register User Request"
// @Success      201      {object}  authDTO.RegisterUserResponse
// @Failure      400      {object}  utils.Response[any] "Error response"
// @Failure      409      {object}  utils.Response[any] "Error response"
// @Failure      500      {object}  utils.Response[any] "Error response"
// @Router       /v1/auth/register [post]
func (h *Handler) RegisterUserHandler(c *gin.Context) error {
	var req authDTO.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Invalid request body",
		})
		return utils.MakeError(errorUc.BadRequest, "invalid request body")
	}

	user, err := h.usecase.RegisterUser(c, req)
	if err != nil {
		// biarkan error dari usecase naik apa adanya agar status code (409/500) tetap sesuai
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to register user",
		})
		return err
	}

	fullName := strings.TrimSpace(user.FirstName + " " + user.LastName)

	h.log.Info(log.LogData{
		Description: "User registered successfully",
		Response:    user,
	})

	c.JSON(http.StatusCreated, authDTO.RegisterUserResponse{
		Data: authDTO.UserResponse{
			ID:    user.ID,
			Name:  fullName,
			Email: user.Email,
		},
	})
	return nil
}

// LoginUserHandler godoc
// @Summary      Login User
// @Description  Login user with email and password
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      authDTO.LoginUserRequest  true  "Login User Request"
// @Success      200      {object}  authDTO.LoginUserResponse
// @Failure      400      {object}  utils.Response[any] "Error response"
// @Failure	  	 401      {object} 	utils.Response[any] "Error response"
// @Failure      500      {object}  utils.Response[any] "Error response"
// @Router       /v1/auth/login [post]
func (h *Handler) LoginUserHandler(c *gin.Context) error {
	var req authDTO.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Invalid request body",
		})
		return utils.MakeError(errorUc.BadRequest, "invalid request body")
	}
	user, err := h.usecase.LoginUser(c, req.Email, req.Password)
	if err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to login user",
		})
		return utils.MakeError(errorUc.InternalServerError, err.Error())
	}

	FullName := user.FirstName + " " + user.LastName
	role := user.Role

	token, err := h.usecase.GenerateToken(user.ID, FullName, role)
	if err != nil {
		h.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to generate token",
		})
		return utils.MakeError(errorUc.InternalServerError, err.Error())
	}
	h.log.Info(log.LogData{
		Description: "User logged in successfully",
		Response:    user,
	})
	c.JSON(http.StatusOK, authDTO.LoginUserResponse{
		User: authDTO.UserResponse{
			ID:    user.ID,
			Name:  FullName,
			Email: user.Email,
		},
		Token: token,
	})
	return nil
}
