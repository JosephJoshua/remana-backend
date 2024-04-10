package logger

import (
	"io"
	stdlog "log"
	"os"
	"runtime/debug"
	"time"

	"github.com/JosephJoshua/repair-management-backend/internal/shared"
	"github.com/rs/zerolog"
)

type loggerError string

func (err loggerError) Error() string {
	return string(err)
}

const (
	ErrLoggerNotInitialized = loggerError("logger not initialized")
)

var log *zerolog.Logger

func Init(logLevel zerolog.Level, appEnv shared.AppEnv) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var output io.Writer

	if appEnv == shared.AppEnvProduction {
		output = os.Stderr
	} else {
		output = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.FieldsExclude = []string{
				"user_agent",
				"git_revision",
				"git_revision_time",
				"dirty_build",
				"go_version",
			}
		})
	}

	var (
		goVersion       string
		gitRevision     string
		gitRevisionTime time.Time
		dirtyBuild      bool
	)

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		goVersion = buildInfo.GoVersion

		for _, v := range buildInfo.Settings {
			switch v.Key {
			case "vcs.revision":
				gitRevision = v.Value

			case "vcs.time":
				t, err := time.Parse(time.RFC3339, v.Value)
				if err != nil {
					continue
				}

				gitRevisionTime = t

			case "vcs.modified":
				dirtyBuild = v.Value == "true"
			}
		}
	}

	l := zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Str("git_revision", gitRevision).
		Time("git_revision_time", gitRevisionTime).
		Bool("dirty_build", dirtyBuild).
		Str("go_version", goVersion).
		Logger()

	log = &l
}

func Get() (zerolog.Logger, error) {
	if log == nil {
		return zerolog.Logger{}, ErrLoggerNotInitialized
	}

	return *log, nil
}

func MustGet() zerolog.Logger {
	l, err := Get()
	if err != nil {
		stdlog.Fatalf("logger.MustGet(); error getting logger: %s", err)
	}

	return l
}
