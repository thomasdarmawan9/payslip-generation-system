package auth

import (
	"context"
	"payslip-generation-system/config/infra"
	"payslip-generation-system/internal/model"
)

type IAuthRepo interface {
	FindByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type AuthRepo struct {
	Infra *infra.Infra
}

func ProvideAuthRepo(infra *infra.Infra) IAuthRepo {
	return &AuthRepo{
		Infra: infra,
	}
}
