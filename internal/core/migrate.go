package core

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/JosephJoshua/remana-backend/internal/shared/logger"
	"github.com/JosephJoshua/remana-backend/internal/shared/projectpath"
	migrate "github.com/rubenv/sql-migrate"
)

const migrationFolder = "db/migrations"

func Migrate(db *sql.DB, dialect string) error {
	migrationSource := &migrate.FileMigrationSource{
		Dir: filepath.Join(projectpath.Root(), migrationFolder),
	}

	l := logger.MustGet()
	l.Info().Msgf("running migrations in %s...", migrationSource.Dir)

	n, err := migrate.Exec(db, dialect, migrationSource, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	l.Info().Int("count", n).Msg("applied migrations")
	return nil
}
