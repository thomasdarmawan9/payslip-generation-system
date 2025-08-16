package infra

import (
	"payslip-generation-system/config"
	"payslip-generation-system/internal/model"
	"payslip-generation-system/pkg/log"

	"gorm.io/gorm"
)

type Infra struct {
	DB *gorm.DB
}

// gunakan provider Postgres yang sudah kita buat sebelumnya
func ProvideInfra(cfg *config.Config, logger *log.LogCustom) *Infra {
	db, err := ProvideDbPostgres(cfg)
	if err != nil {
		logger.Error(log.LogData{
			Err:         err,
			Description: "failed to connect to PostgreSQL database",
			StartTime:   nil,
			Response:    nil,
		})
		panic("cannot start app without DB")
	}

	infra := &Infra{
		DB: db,
	}

	if cfg.DBConfig.EnableAutoMigration {
		// taruh semua migrasi model di sini
		if err := infra.DB.AutoMigrate(&model.User{}); err != nil {
			logger.Error(log.LogData{
				Err:         err,
				Description: "database migration failed",
				StartTime:   nil,
				Response:    nil,
			})
			panic("auto migration failed")
		}
		logger.Info(log.LogData{
			Description: "database migration completed successfully",
			StartTime:   nil,
			Response:    nil,
		})
	}

	return infra
}
