package infra

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"backend/internal/model"
)

func NewPostgresConnection(lc fx.Lifecycle, logger *Logger, cfg *Config) (*gorm.DB, error) {
	var gormConfig *gorm.Config
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      gormLogger.Info,
			Colorful:      true,
		},
	)
	gormConfig = &gorm.Config{
		TranslateError: true,
		Logger:         newLogger,
	}
	db, err := gorm.Open(postgres.Open(cfg.DbUrl), gormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running migrations with db url: " + cfg.DbUrl)

			// run migrations

			err = db.AutoMigrate(&model.User{})

			logger.Info("migrations applied")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			dbsql, err := db.DB()
			if err != nil {
				logger.Error()
			}

			err = dbsql.Close()
			if err != nil {
				logger.Error()
			}

			logger.Info("db connection closed")

			return nil
		},
	})

	return db, nil
}
