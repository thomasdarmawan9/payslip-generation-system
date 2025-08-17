package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	authmidware "payslip-generation-system/internal/middleware"
)

func (r *Route) SetupRoute(router *gin.Engine) {
	// Health check
	router.GET("/health-check", healthCheck)

	// CORS
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = r.Cfg.Cors.AllowOrigins
	configCors.AllowHeaders = r.Cfg.Cors.AllowHeaders
	configCors.AllowMethods = r.Cfg.Cors.AllowMethods
	configCors.AllowCredentials = r.Cfg.Cors.AllowCredentials
	router.Use(cors.New(configCors))

	// V1
	v1 := router.Group("/v1")

	// Public auth
	auth := v1.Group("/auth")
	auth.POST("/register", r.processTimeout(WrapWithErrorHandler(r.handler.RegisterUserHandler), 5*time.Second))
	auth.POST("/login", r.processTimeout(WrapWithErrorHandler(r.handler.LoginUserHandler), 5*time.Second))

	// Protected (JWT) — apply middleware.Auth
	protected := v1.Group("")
	authmidware.New(protected, r.Cfg, r.Log) // ini memasang AuthJwt untuk semua route di bawahnya

	// ADMIN only group
	admin := protected.Group("")
	admin.Use(RequireAdmin()) // helper kecil di bawah
	// contoh endpoint admin (buat period payroll)
	admin.POST("/payroll/periods", r.processTimeout(WrapWithErrorHandler(r.handler.CreateAttendancePeriodHandler), 10*time.Second))
	admin.POST("/payroll/periods/:period_id/run", r.processTimeout(WrapWithErrorHandler(r.handler.RunPayrollHandler), 30*time.Second))
	// USER or ADMIN
	user := protected.Group("")
	user.Use(RequireUserOrAdmin())
	// contoh endpoint submit attendance
	user.POST("/attendance/submit", r.processTimeout(WrapWithErrorHandler(r.handler.SubmitAttendanceHandler), 10*time.Second))
	user.POST("/overtime/submit", r.processTimeout(WrapWithErrorHandler(r.handler.SubmitOvertimeHandler), 10*time.Second))
	user.POST("/reimbursements", r.processTimeout(WrapWithErrorHandler(r.handler.CreateReimbursementHandler), 10*time.Second))
	user.GET("/payslips/periods/:period_id",
		r.processTimeout(WrapWithErrorHandler(r.handler.GeneratePayslipHandler), 10*time.Second))

}

// Health
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP", "message": "Service is running"})
}

// Error wrapper
func WrapWithErrorHandler(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			// kalau handler sudah nulis response, jangan double-write.
			if !c.IsAborted() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"responseCode":    "5000100",
					"responseMessage": err.Error(),
				})
				c.Abort()
			}
		}
	}
}

// Timeout wrapper — versi lebih aman
func (r *Route) processTimeout(handler gin.HandlerFunc, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		processDone := make(chan struct{})
		go func() {
			handler(c)
			processDone <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"responseCode":    "4080100",
				"responseMessage": "Request Process Timeout",
			})
		case <-processDone:
			// success
		}
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role := c.GetString("role"); role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"responseCode":    "4030100",
				"responseMessage": "admin only",
			})
			return
		}
		c.Next()
	}
}
func RequireUserOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role := c.GetString("role"); role != "user" && role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"responseCode":    "4030101",
				"responseMessage": "forbidden",
			})
			return
		}
		c.Next()
	}
}
