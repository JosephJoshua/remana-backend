package core

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/JosephJoshua/remana-backend/internal/shared/projectpath"
	migrate "github.com/rubenv/sql-migrate"
)

const migrationFolder = "db/migrations"

var migrationSource = &migrate.FileMigrationSource{
	Dir: filepath.Join(projectpath.Root(), migrationFolder),
}

func GetPendingMigrationCount(db *sql.DB, dialect string) (int, error) {
	plannedMigrations, _, err := migrate.PlanMigration(db, dialect, migrationSource, migrate.Up, -1)
	if err != nil {
		return 0, fmt.Errorf("failed to plan migration: %w", err)
	}

	return len(plannedMigrations), nil
}

func Migrate(ctx context.Context, db *sql.DB, dialect string) (int, error) {
	n, err := migrate.ExecContext(ctx, db, dialect, migrationSource, migrate.Up)
	if err != nil {
		return 0, fmt.Errorf("failed to execute migration: %w", err)
	}

	return n, nil
}
