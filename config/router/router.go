package router

import (
	"context"
	"net/http"

	// authmidware "payslip-generation-system/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (r *Route) SetupRoute(router *gin.Engine) {
	// Health check
	router.GET("/health-check", healthCheck)

	// CORS setup
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = r.Cfg.Cors.AllowOrigins
	configCors.AllowHeaders = r.Cfg.Cors.AllowHeaders
	configCors.AllowMethods = r.Cfg.Cors.AllowMethods
	configCors.AllowCredentials = r.Cfg.Cors.AllowCredentials
	router.Use(cors.New(configCors))

	// V1 API group
	v1 := router.Group("/v1")

	// Auth Route group
	auth := v1.Group("/auth")

	auth.POST("/register", r.processTimeout(WrapWithErrorHandler(r.handler.RegisterUserHandler), 5*time.Second))
	auth.POST("/login", r.processTimeout(WrapWithErrorHandler(r.handler.LoginUserHandler), 5*time.Second))

	// Invoice route group
	// invoice := v1.Group("/invoices")

	// Middleware for authentication
	// authmidware.New(invoice, *r.Cfg, r.Log)
}

// healthCheck returns the health status of the service.
// @Summary Health Check
// @Description Returns the health status of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health-check [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"message": "Service is running",
	})
}

func WrapWithErrorHandler(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"responseCode":    "5000100",
				"responseMessage": err.Error(),
			})
		}
	}
}

// processTimeout wraps a gin.HandlerFunc with a context timeout.
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
