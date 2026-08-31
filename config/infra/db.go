package infra

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"payslip-generation-system/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProvideDbPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" && cfg != nil {
		dsn = strings.TrimSpace(cfg.DBConfig.DBPostgresConfig["postgres"])
	}
	if dsn == "" {
		return nil, fmt.Errorf("database URL is not configured")
	}

	startTime := time.Now()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	mode := "unknown"
	if cfg != nil {
		mode = cfg.AppEnvMode.Mode
	}
	logMessage := fmt.Sprintf("connect db with %v", mode)

	if err != nil {
		log.Printf("[ERROR] %s | duration: %v | error: %v\n", logMessage, time.Since(startTime), err)
		return nil, err
	}

	log.Printf("[INFO] %s | duration: %v | connection successful\n", logMessage, time.Since(startTime))
	return db, nil
}
