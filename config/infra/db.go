package infra

import (
	"fmt"
	"log"
	"payslip-generation-system/config"
	"payslip-generation-system/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProvideDbPostgres(cfg *config.Config) (*gorm.DB, error) {
	var dsn string

	switch cfg.AppEnvMode.Mode {
	case utils.DEV, utils.DEV_TEST, utils.PROD:
		dsn = cfg.DBConfig.DBPostgresConfig["postgres"]
	default:
		dsn = cfg.DBConfig.DBPostgresConfig["postgres"]
	}

	startTime := time.Now()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	logMessage := fmt.Sprintf("connect db with %v using %v", cfg.AppEnvMode.Mode, dsn)

	if err != nil {
		log.Printf("[ERROR] %s | duration: %v | error: %v\n", logMessage, time.Since(startTime), err)
		return nil, err
	}

	log.Printf("[INFO] %s | duration: %v | connection successful\n", logMessage, time.Since(startTime))
	return db, nil
}
