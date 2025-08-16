package router

import (
	"payslip-generation-system/config"
	"payslip-generation-system/internal/handler"
	"payslip-generation-system/internal/usecase"
	"payslip-generation-system/pkg/log"
)

type Route struct {
	Cfg     *config.Config
	Log     *log.LogCustom
	handler handler.Handler
	usecase usecase.IUsecase
}

func ProvideRoute(
	cfg *config.Config,
	log *log.LogCustom,
	handler handler.Handler,
	usecase usecase.IUsecase,
) Route {
	return Route{
		Cfg:     cfg,
		Log:     log,
		handler: handler,
		usecase: usecase,
	}
}
