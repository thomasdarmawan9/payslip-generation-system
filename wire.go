//go:build wireinject
// +build wireinject

package main

import (
	"payslip-generation-system/config"
	"payslip-generation-system/config/infra"
	"payslip-generation-system/config/router"
	"payslip-generation-system/internal/handler"
	auth "payslip-generation-system/internal/repository/auth"
	"payslip-generation-system/internal/repository/tx"
	"payslip-generation-system/internal/usecase"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/transport"

	"github.com/google/wire"
)

var Configs = wire.NewSet(
	config.ProvideConfig,
)

var LoggerSet = wire.NewSet(
	log.ProvideLogger,
)

var InfraSet = wire.NewSet(
	infra.ProvideInfra,
)

var RepoSet = wire.NewSet(
	auth.ProvideAuthRepo,
	tx.ProvideTxManager,
)

var InternalDomain = wire.NewSet(
	RepoSet,
	usecase.ProvideUsc,
)

var Handler = wire.NewSet(
	handler.ProvideHandler,
)

var Server = wire.NewSet(
	Handler,
	router.ProvideRoute,
	transport.ProvideHttp,
)

func ServerApp() *transport.HTTP {
	wire.Build(
		Configs,
		InfraSet,
		LoggerSet,
		InternalDomain,
		Server)

	return &transport.HTTP{}
}
