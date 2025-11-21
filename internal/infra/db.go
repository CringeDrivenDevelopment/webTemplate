package infra

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"

	projectroot "backend"
)

func NewPostgresConnection(lc fx.Lifecycle, logger *Logger, cfg *Config) (*pgxpool.Pool, error) {
	ctxWithCancel, cancel := context.WithCancel(context.Background())

	pool, err := pgxpool.New(ctxWithCancel, cfg.DbUrl)
	if err != nil {
		cancel()
		return nil, err
	}

	// configure pool
	poolConfig := pool.Config()
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30
	poolConfig.HealthCheckPeriod = time.Minute

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// check if online
			if err := pool.Ping(ctx); err != nil {
				return err
			}

			// run migrations
			goose.SetBaseFS(projectroot.EmbedMigrations)
			goose.SetLogger(&ZapGooseAdapter{zap: logger.Zap})
			if err := goose.SetDialect("postgres"); err != nil {
				return err
			}
			db := stdlib.OpenDBFromPool(pool)
			if err := goose.Up(db, "sql/migrations"); err != nil {
				return err
			}
			if err := db.Close(); err != nil {
				return err
			}

			logger.Info("migrations applied")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			cancel()

			logger.Info("db connection closed")

			return nil
		},
	})

	return pool, nil
}
