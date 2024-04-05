package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephJoshua/repair-management-backend/internal/genapi"
	"github.com/JosephJoshua/repair-management-backend/internal/logger"
	"github.com/JosephJoshua/repair-management-backend/internal/shared"
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

	oasServer, err := genapi.NewServer(shared.Server{}, nil)
	if err != nil {
		return fmt.Errorf("error creating server: %w", err)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           oasServer,
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

type appConfig struct {
	ServerAddr string `mapstructure:"remana_server_addr"`
}

func loadConfig() (appConfig, error) {
	viper.SetConfigFile(".env")

	viper.SetDefault("remana_server_addr", "localhost:8080")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return appConfig{}, fmt.Errorf("error reading in config: %w", err)
	}

	var config appConfig
	if err = viper.Unmarshal(&config); err != nil {
		return appConfig{}, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return config, nil
}

func main() {
	logger.Init(zerolog.DebugLevel)
	log := logger.MustGet()

	config, err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	ctx := context.Background()
	if err = run(ctx, config.ServerAddr); err != nil {
		log.Fatal().Err(err).Msg("error running app")
	}
}
