package usecase

import (
	"payslip-generation-system/config"
	"payslip-generation-system/internal/model"
	repositoryAuth "payslip-generation-system/internal/repository/auth"
	"payslip-generation-system/internal/repository/tx"
	"payslip-generation-system/pkg/log"

	authDTO "payslip-generation-system/internal/dto/auth"

	"github.com/gin-gonic/gin"
)

type IUsecase interface {
	LoginUser(ctx *gin.Context, email, password string) (*model.User, error)
	GenerateToken(userID uint, name string) (string, error)

	RegisterUser(ctx *gin.Context, userDTO authDTO.RegisterUserRequest) (*model.User, error)
}

type usecase struct {
	log       *log.LogCustom
	cfg       *config.Config
	txManager tx.TxManager
	authRepo  repositoryAuth.IAuthRepo
}

func ProvideUsc(log *log.LogCustom, cfg *config.Config,
	txManager tx.TxManager,
	authRepo repositoryAuth.IAuthRepo,
) IUsecase {
	return &usecase{
		log:       log,
		cfg:       cfg,
		txManager: txManager,
		authRepo:  authRepo,
	}
}
