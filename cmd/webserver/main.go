package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/fs"
	stdlog "log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/infrastructure/core"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/projectpath"
	"github.com/JosephJoshua/remana-backend/internal/shared"
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

func Run(ctx context.Context, db *pgxpool.Pool, addr string, certPEM string, keyPEM string) error {
	log := logger.MustGet()

	srv, middlewares, err := core.NewAPIServer(db)
	if err != nil {
		return fmt.Errorf("error creating server: %w", err)
	}

	handler := http.Handler(srv)

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return fmt.Errorf("error loading cert and key: %w", err)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: ReadHeaderTimeout,
		IdleTimeout:       IdleTimeout,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		TLSConfig: &tls.Config{
			Certificates:             []tls.Certificate{cert},
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	listenErr := make(chan error)

	go func() {
		log.Info().Msgf("server listening on %s", httpServer.Addr)

		if err = httpServer.ListenAndServeTLS("", ""); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
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

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing conn string: %w", err)
	}

	config.ConnConfig.Tracer = &logger.PgxLogTracer{}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return conn, nil
}

type appConfig struct {
	ServerAddr   string        `mapstructure:"remana_server_addr"    validate:"required"`
	AppEnv       shared.AppEnv `mapstructure:"remana_app_env"        validate:"required"`
	ConnString   string        `mapstructure:"remana_conn_string"    validate:"required"`
	CertFilePath string        `mapstructure:"remana_cert_file_path" validate:"required"`
	KeyFilePath  string        `mapstructure:"remana_key_file_path"  validate:"required"`
}

func loadConfig() (appConfig, error) {
	viper.SetConfigFile(".env")

	viper.SetDefault("remana_server_addr", "localhost:8080")
	viper.SetDefault("remana_app_env", "production")
	viper.SetDefault("remana_cert_file_path", "server.crt")
	viper.SetDefault("remana_key_file_path", "server.key")

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

	l.Info().Str("mode", string(config.AppEnv)).Msgf("app running in %s mode", config.AppEnv)

	ctx := context.Background()

	pool, err := connectDB(ctx, config.ConnString)
	if err != nil {
		l.Fatal().Err(err).Msg("error connecting to database")
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)
	n, err := core.GetPendingMigrationCount(db, "postgres")

	if err != nil {
		l.Panic().Err(err).Msg("error checking if migration is needed")
	}

	if n > 0 {
		l.Warn().Int("count", n).Msg("there are pending migrations")
	}

	certFilePath, err := url.JoinPath(projectpath.Root(), config.CertFilePath)
	if err != nil {
		l.Panic().Err(err).Msg("error joining cert file path")
	}

	keyFilePath, err := url.JoinPath(projectpath.Root(), config.KeyFilePath)
	if err != nil {
		l.Panic().Err(err).Msg("error joining key file path")
	}

	certPEM, err := os.ReadFile(certFilePath)
	if err != nil {
		l.Panic().Err(err).Msg("error reading cert file")
	}

	keyPEM, err := os.ReadFile(keyFilePath)
	if err != nil {
		l.Panic().Err(err).Msg("error reading key file")
	}

	if err = Run(ctx, pool, config.ServerAddr, string(certPEM), string(keyPEM)); err != nil {
		l.Panic().Err(err).Msg("error running app")
	}
}
