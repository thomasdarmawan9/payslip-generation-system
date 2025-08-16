package usecase

import (
	"payslip-generation-system/config"
	"payslip-generation-system/internal/model"
	atRepo "payslip-generation-system/internal/repository/attendance"
	apRepo "payslip-generation-system/internal/repository/attendanceperiod"
	otRepo "payslip-generation-system/internal/repository/overtime"
	rbRepo "payslip-generation-system/internal/repository/reimbursement"
	repositoryAuth "payslip-generation-system/internal/repository/auth"
	repoTx "payslip-generation-system/internal/repository/tx"
	"payslip-generation-system/pkg/log"

	authDTO "payslip-generation-system/internal/dto/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IUsecase interface {
	LoginUser(ctx *gin.Context, email, password string) (*model.User, error)
	GenerateToken(userID uint, name, role string) (string, error)
	RegisterUser(ctx *gin.Context, userDTO authDTO.RegisterUserRequest) (*model.User, error)

	CreateAttendancePeriod(ctx *gin.Context, name, start, end string) (*model.AttendancePeriod, error)
	SubmitAttendance(ctx *gin.Context, userID uint, dateStr string) (*model.Attendance, bool, error)

	SubmitOvertime(ctx *gin.Context, userID uint, dateStr string, hours float64) (*model.Overtime, bool, error)
	CreateReimbursement(ctx *gin.Context, userID uint, dateStr string, amount float64, description string) (*model.Reimbursement, error)
}

type usecase struct {
	cfg       *config.Config
	log       *log.LogCustom
	authRepo  repositoryAuth.IAuthRepo
	txManager repoTx.TxManager
	apRepo    apRepo.Repo
	atRepo    atRepo.Repo
	otRepo otRepo.Repo
	rbRepo rbRepo.Repo
}

func ProvideUsc(
	cfg *config.Config,
	l *log.LogCustom,
	db *gorm.DB,
	authRepo repositoryAuth.IAuthRepo,
	txManager repoTx.TxManager,
) IUsecase {
	u := &usecase{
		cfg:       cfg,
		log:       l,
		authRepo:  authRepo,
		txManager: txManager,
	}
	// inject attendance repos
	u.apRepo = apRepo.New(db)
	u.atRepo = atRepo.New(db)
	u.otRepo = otRepo.New(db)
	u.rbRepo = rbRepo.New(db)
	return u
}
