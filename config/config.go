package config

import (
	"payslip-generation-system/pkg/dbconfig"
	"payslip-generation-system/pkg/env"
	"payslip-generation-system/pkg/log"
	"time"
)

type Config struct {
	EnvConfig  env.Config      `mapstructure:"envLib"`
	AppEnvMode AppEnvMode      `mapstructure:"appEnvMode"`
	DBConfig   dbconfig.Config `mapstructure:"databaseConfig"`
	LogConfig  log.Config      `mapstructure:"logConfig"`
	ConfigEnv  env.EnvConfig

	Server struct {
		Shutdown struct {
			CleanupPeriodSeconds int64 `mapstructure:"cleanup_period_seconds"`
			GracePeriodSeconds   int64 `mapstructure:"grace_period_seconds"`
		} `mapstructure:"shutdown"`
		Timeout struct {
			Duration time.Duration `mapstructure:"duration"`
		}
	} `mapstructure:"server"`

	Cors CORSConfig `mapstructure:"cors"`
}

type AppEnvMode struct {
	Mode           string `mapstructure:"mode"`        // dev, prod, dev_test
	GinMode        string `mapstructure:"ginMode"`     // debug, release
	DebugMode      bool   `mapstructure:"debugMode"`   // true, false
	IsPrettyLog    bool   `mapstructure:"isPrettyLog"` // true, false
	TestPathPrefix string `mapstructure:"testPathPrefix"`
	Port           string `mapstructure:"port"` // e.g., 9898
	Host           string `mapstructure:"host"` // e.g., "localhost"
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allowOrigins"`
	AllowMethods     []string `mapstructure:"allowMethods"`
	AllowHeaders     []string `mapstructure:"allowHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	ExposeHeaders    []string `mapstructure:"exposeHeaders"`
}
