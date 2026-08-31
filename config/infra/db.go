package infra

import (
	"fmt"
	"log"
	"os"
	"payslip-generation-system/config"
	"payslip-generation-system/utils"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProvideDbPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" && cfg != nil {
		dsn = cfg.DBConfig.DBPostgresConfig["postgres"]
	}
	if dsn == "" {
		return nil, fmt.Errorf("database URL is not configured")
	}

	startTime := time.Now()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	logMessage := fmt.Sprintf("connect db with %v", cfg.AppEnvMode.Mode)

	if err != nil {
		log.Printf("[ERROR] %s | duration: %v | error: %v\n", logMessage, time.Since(startTime), err)
		return nil, err
	}

	log.Printf("[INFO] %s | duration: %v | connection successful\n", logMessage, time.Since(startTime))
	return db, nil
}
