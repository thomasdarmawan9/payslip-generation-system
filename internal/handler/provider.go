package handler

import (
	"payslip-generation-system/config"
	"payslip-generation-system/internal/usecase"
	"payslip-generation-system/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     *config.Config
	log     *log.LogCustom
	usecase usecase.IUsecase
	ctx     *gin.Context
}

func ProvideHandler(
	cfg *config.Config,
	l *log.LogCustom,
	usecase usecase.IUsecase,
) Handler {
	return Handler{
		cfg:     cfg,
		log:     l,
		usecase: usecase,
	}
}
