package transport

import (
	"fmt"
	"os"
	"os/signal"
	"payslip-generation-system/config"
	"payslip-generation-system/config/router"
	invoiceLog "payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"payslip-generation-system/docs"
)

// ServerState is an indicator if this server's state.
type ServerState int

const (
	// ServerStateReady indicates that the server is ready to serve.
	ServerStateReady ServerState = iota + 1
	// ServerStateInGracePeriod indicates that the server is in its grace
	// period and will shut down after it is done cleaning up.
	ServerStateInGracePeriod
	// ServerStateInCleanupPeriod indicates that the server no longer
	// responds to any requests, is cleaning up its internal state, and
	// will shut down shortly.
	ServerStateInCleanupPeriod
)

// HTTP is the HTTP server.
type HTTP struct {
	Config *config.Config
	Route  router.Route
	State  ServerState
	Server *gin.Engine
	Log    *invoiceLog.LogCustom
}

func ProvideHttp(Config *config.Config,
	route router.Route,
	log *invoiceLog.LogCustom,
) *HTTP {
	srv := gin.New()
	switch Config.AppEnvMode.Mode {
	case utils.DEV, utils.DEV_TEST:
		gin.SetMode(gin.DebugMode)
	case utils.PROD:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	srv.Use(gin.Logger())
	srv.Use(gin.Recovery())

	srv.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"responseCode": "40400000", "responseMessage": "Invalid Path"})
	})

	return &HTTP{
		Config: Config,
		Route:  route,
		Server: srv,
		Log:    log,
	}
}

func (h *HTTP) Serve() {
	h.Route.SetupRoute(h.Server)
	h.setupGracefulShutdown()
	h.State = ServerStateReady

	addr := h.Config.AppEnvMode.Host + ":" + h.Config.AppEnvMode.Port
	h.setupSwaggerDocs()

	err := h.Server.Run(addr)
	if err != nil {
		log.Fatal().Msg("Failed to start server")
	}
}

func (h *HTTP) setupGracefulShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go h.respondToSigterm(done)
}

func (h *HTTP) respondToSigterm(done chan os.Signal) {
	<-done
	defer os.Exit(0)

	shutdownConfig := h.Config.Server.Shutdown

	log.Info().Msg("Received SIGTERM.")
	log.Info().Int64("seconds", shutdownConfig.GracePeriodSeconds).Msg("Entering grace period.")
	h.State = ServerStateInGracePeriod
	time.Sleep(time.Duration(shutdownConfig.GracePeriodSeconds) * time.Second)

	log.Info().Int64("seconds", shutdownConfig.CleanupPeriodSeconds).Msg("Entering cleanup period.")
	h.State = ServerStateInCleanupPeriod
	time.Sleep(time.Duration(shutdownConfig.CleanupPeriodSeconds) * time.Second)

	log.Info().Msg("Cleaning up completed. Shutting down now.")
}

func (h *HTTP) setupSwaggerDocs() {
	if h.Config.AppEnvMode.Mode == utils.DEV || h.Config.AppEnvMode.Mode == utils.DEV_TEST {
		docs.SwaggerInfo.Title = h.Config.EnvConfig.AppConfig.Name
		docs.SwaggerInfo.Version = h.Config.EnvConfig.AppConfig.Version

		host := h.Config.AppEnvMode.Host
		port := h.Config.AppEnvMode.Port

		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "9898"
		}

		swaggerDocURL := fmt.Sprintf("http://%s:%s/swagger/doc.json", host, port)

		swaggerOpt := ginSwagger.URL(swaggerDocURL)

		h.Server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerOpt))

		log.Info().Str("url", swaggerDocURL).Msg("Swagger documentation enabled.")
	}
}
