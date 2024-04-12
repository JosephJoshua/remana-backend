package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/core"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/JosephJoshua/remana-backend/internal/shared/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	ReadHeaderTimeout = 5 * time.Second
	IdleTimeout       = 60 * time.Second
	ReadTimeout       = 30 * time.Second
	WriteTimeout      = 30 * time.Second
	ShutdownTimeout   = 10 * time.Second
)

func run(ctx context.Context, addr string) error {
	log := logger.MustGet()

	srv, middlewares, err := core.NewAPIServer()
	if err != nil {
		return fmt.Errorf("error creating server: %w", err)
	}

	handler := http.Handler(srv)

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: ReadHeaderTimeout,
		IdleTimeout:       IdleTimeout,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
	}

	listenErr := make(chan error)

	go func() {
		log.Info().Msgf("server listening on %s", httpServer.Addr)

		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenErr <- fmt.Errorf("error listening and serving: %w", err)
		}
	}()

	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err = <-listenErr:
		return err

	case <-signalCtx.Done():
		log.Info().Msg("server shutting down")

		shutdownCtx, cancel := context.WithTimeout(ctx, ShutdownTimeout)
		defer cancel()

		if err = httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down http server: %w", err)
		}

		log.Info().Msg("server shut down")
	}

	return nil
}

func connectDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	log := logger.MustGet()

	log.Info().Msg("connecting to database...")

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return conn, nil
}

type appConfig struct {
	ServerAddr string        `mapstructure:"remana_server_addr" validate:"required"`
	AppEnv     shared.AppEnv `mapstructure:"remana_app_env"     validate:"required"`
	ConnString string        `mapstructure:"remana_conn_string" validate:"required"`
}

func loadConfig() (appConfig, error) {
	viper.SetConfigFile(".env")

	viper.SetDefault("remana_server_addr", "localhost:8080")
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
	log := logger.MustGet()

	log.Info().Str("mode", string(config.AppEnv)).Msgf("app running in %s mode", config.AppEnv)

	ctx := context.Background()

	pool, err := connectDB(ctx, config.ConnString)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)
	n, err := core.GetPendingMigrationCount(db, "postgres")

	if err != nil {
		log.Panic().Err(err).Msg("error checking if migration is needed")
	}

	if n > 0 {
		log.Warn().Int("count", n).Msg("there are pending migrations")
	}

	if err = run(ctx, config.ServerAddr); err != nil {
		log.Panic().Err(err).Msg("error running app")
	}
}
