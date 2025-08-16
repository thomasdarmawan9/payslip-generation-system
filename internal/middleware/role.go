package middleware

import (
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrForbidden, "admin only"))))
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireUserOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "user" && role != "admin" {
			utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrForbidden))))
			c.Abort()
			return
		}
		c.Next()
	}
}
