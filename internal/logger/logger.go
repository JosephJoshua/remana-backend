package logger

import (
	stdlog "log"
	"runtime/debug"

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

func Init(logLevel zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	output := zerolog.NewConsoleWriter()

	var (
		goVersion   string
		gitRevision string
		vcsModified bool
	)

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		goVersion = buildInfo.GoVersion

		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.version" {
				gitRevision = v.Value
			} else if v.Key == "vcs.modified" {
				vcsModified = v.Value == "true"
			}
		}
	}

	l := zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Str("git_revision", gitRevision).
		Bool("vcs_modified", vcsModified).
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
		stdlog.Fatalf("logger.MustGet() > error getting logger: %s", err)
	}

	return l
}
