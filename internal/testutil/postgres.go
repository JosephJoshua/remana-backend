package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/infrastructure/core"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func StartPostgresContainer(pool *dockertest.Pool) (*dockertest.Resource, *pgxpool.Pool, error) {
	const (
		pgUsername = "username"
		pgPassword = "secretpassword"
		pgDBName   = "remana"

		pgContainerLifetimeSecs = 30
		pgContainerMaxWait      = 30 * time.Second
	)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", pgUsername),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPassword),
			fmt.Sprintf("POSTGRES_DB=%s", pgDBName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error running postgres container: %w", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUsername, pgPassword, hostAndPort, pgDBName)

	if err = resource.Expire(pgContainerLifetimeSecs); err != nil {
		return nil, nil, fmt.Errorf("error setting expiry date on postgres container: %w", err)
	}

	pool.MaxWait = pgContainerMaxWait

	var db *pgxpool.Pool

	if err = pool.Retry(func() error {
		config, configErr := pgxpool.ParseConfig(databaseURL)
		if configErr != nil {
			return fmt.Errorf("error parsing database URL: %w", err)
		}

		if _, ok := logger.Get(); ok {
			config.ConnConfig.Tracer = &logger.PgxLogTracer{}
		}

		db, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			return err
		}

		return db.Ping(context.Background())
	}); err != nil {
		return nil, nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return resource, db, nil
}

func MigratePostgres(ctx context.Context, db *pgxpool.Pool) error {
	const maxWait = 5 * time.Second

	rawDB := stdlib.OpenDBFromPool(db)

	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	if _, err := core.Migrate(ctx, rawDB, "postgres"); err != nil {
		return err
	}

	return nil
}
