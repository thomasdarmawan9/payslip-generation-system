package auth

import (
	"context"
	"errors"
	"payslip-generation-system/internal/model"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"

	repotx "payslip-generation-system/internal/repository/tx"
)

func (r *AuthRepo) FindByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	if err := r.Infra.DB.WithContext(ctx).Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}
	return &user, nil
}

func (r *AuthRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.Infra.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}
	return &user, nil
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *model.User) error {
	db := repotx.GetDB(ctx, r.Infra.DB)
	if err := db.Create(user).Error; err != nil {
		// Unique constraint violation
		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

var ErrEmailAlreadyExists = errors.New("email already registered")

func isUniqueViolation(err error) bool {
	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		return pgxErr.Code == "23505"
	}
	if pqErr, ok := err.(*pq.Error); ok {
		return string(pqErr.Code) == "23505"
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
