package usecase

import (
	"fmt"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"
	"strings"
	"time"

	errorUc "payslip-generation-system/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"

	authDTO "payslip-generation-system/internal/dto/auth"
)

// RegisterUser membuat akun baru lalu mengembalikan entity user (transactional).
func (u *usecase) RegisterUser(ctx *gin.Context, req authDTO.RegisterUserRequest) (*model.User, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// --- BEGIN TX ---
	txCtx, err := u.txManager.Begin(ctx)
	if err != nil {
		u.log.Error(log.LogData{Err: err, Description: "failed to begin transaction"})
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to begin transaction")
	}
	defer func() {
		// jika ada error di return path, rollback; kalau tidak, commit
		if err != nil {
			_ = u.txManager.Rollback(txCtx)
		} else {
			if cmErr := u.txManager.Commit(txCtx); cmErr != nil {
				_ = u.txManager.Rollback(txCtx)
				err = utils.MakeError(errorUc.InternalServerError, "failed to commit transaction")
			}
		}
	}()
	// --- END TX setup ---

	// Cek existing user by email
	existing, err := u.authRepo.FindByEmail(txCtx, email)
	if err != nil {
		u.log.Error(log.LogData{Err: err, Description: "failed to check existing user by email"})
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to check existing user")
	}
	if existing != nil {
		u.log.Info(log.LogData{Description: "email already registered", Response: email})
		return nil, utils.MakeError(errorUc.ConflictError, "email already registered")
	}

	// Hash password
	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		u.log.Error(log.LogData{Err: hashErr, Description: "failed to hash password"})
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to hash password")
	}

	// Tentukan profile complete
	isComplete := req.IsProfileComplete
	if !isComplete {
		isComplete = email != "" && req.FirstName != "" && req.LastName != ""
	}

	user := &model.User{
		Email:             email,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		ProfileImageURL:   req.ProfileImageURL,
		PasswordHash:      string(hashed),
		GoogleID:          req.GoogleID,
		Age:               req.Age,
		Bio:               req.Bio,
		Location:          req.Location,
		Interests:         req.Interests, // pq.StringArray -> text[]
		Role:              req.Role,
		Salary:            req.Salary,
		IsProfileComplete: isComplete,
		// CreatedAt/UpdatedAt by GORM
	}

	// Simpan user (masih dalam tx)
	if err = u.authRepo.CreateUser(txCtx, user); err != nil {
		u.log.Error(log.LogData{Err: err, Description: "failed to create user"})
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "duplicate key") || strings.Contains(low, "unique constraint") {
			return nil, utils.MakeError(errorUc.ConflictError, "email already registered")
		}
		return nil, utils.MakeError(errorUc.InternalServerError, "failed to create user")
	}

	u.log.Info(log.LogData{Description: "user registered successfully", Response: user})
	return user, nil
}

func (u *usecase) LoginUser(ctx *gin.Context, email, password string) (*model.User, error) {

	user, err := u.authRepo.FindByEmail(ctx, email)
	if err != nil {
		u.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to retrieve user",
		})
		return nil, utils.MakeError(errorUc.InternalServerError, err.Error())
	}
	if user == nil {
		u.log.Error(log.LogData{
			Description: "User not found",
		})
		return nil, utils.MakeError(errorUc.NotFoundError, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		u.log.Error(log.LogData{
			Description: "Invalid password",
			Err:         err,
		})
		return nil, utils.MakeError(errorUc.ErrUnauthorized, "invalid email or password")
	}

	u.log.Info(log.LogData{
		Description: fmt.Sprintf("User %s logged in successfully", user.Email),
		Response:    user,
	})

	return user, nil
}

func (u *usecase) GenerateToken(userID uint, name, role string) (string, error) {
	token, err := GenerateToken(userID, name, role)
	if err != nil {
		u.log.Error(log.LogData{
			Err:         err,
			Description: "Failed to generate token",
		})
		return "", utils.MakeError(errorUc.InternalServerError, err.Error())
	}
	u.log.Info(log.LogData{
		Description: "Token generated successfully",
		Response:    token,
	})
	return token, nil
}

var jwtSecret = []byte("A7M+TXRMxdz0N3nFLjGaxVKgkELowtbxWipS+IFZkVE=") // Ganti dengan env di production

func GenerateToken(userID uint, name, role string) (string, error) {
	// Define token expiration (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"role":    role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
