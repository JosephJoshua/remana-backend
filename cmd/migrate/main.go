package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	stdlog "log"

	"github.com/JosephJoshua/remana-backend/internal/core"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/JosephJoshua/remana-backend/internal/shared/logger"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type appConfig struct {
	AppEnv     shared.AppEnv `mapstructure:"remana_app_env"     validate:"required"`
	ConnString string        `mapstructure:"remana_conn_string" validate:"required"`
}

func connectDB(connString string) (*sql.DB, error) {
	log := logger.MustGet()

	log.Info().Msg("connecting to database...")

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

func loadConfig() (appConfig, error) {
	viper.SetConfigFile(".env")
	viper.SetDefault("remana_app_env", "production")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return appConfig{}, fmt.Errorf("error reading in config: %w", err)
	}

	var config appConfig
	if err = viper.Unmarshal(&config); err != nil {
		return appConfig{}, fmt.Errorf("error unmarshalling config: %w", err)
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return appConfig{}, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		stdlog.Fatalf("error loading config: %v", err)
	}

	logger.Init(zerolog.DebugLevel, config.AppEnv)
	l := logger.MustGet()

	ctx := context.Background()

	db, err := connectDB(config.ConnString)
	if err != nil {
		l.Fatal().Err(err).Msg("error connecting to database")
	}
	defer func() {
		if err = db.Close(); err != nil {
			l.Panic().Err(err).Msg("error closing database connection")
		}
	}()

	if n, migrateErr := core.Migrate(ctx, db, "postgres"); migrateErr != nil {
		l.Panic().Err(migrateErr).Msg("error migrating database")
	} else {
		l.Info().Int("count", n).Msg("applied migrations")
	}
}
