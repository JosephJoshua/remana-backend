package core

import (
	"context"
	"net/http"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

type correlationIDKey struct{}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := logger.MustGet()
		correlationID := xid.New().String()

		ctx := context.WithValue(r.Context(), correlationIDKey{}, correlationID)

		r = r.WithContext(ctx)
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("correlation_id", correlationID)
		})

		w.Header().Add("X-Correlation-ID", correlationID)

		lrw := newLoggingResponseWriter(w)
		r = r.WithContext(l.WithContext(r.Context()))

		defer func() {
			err := recover()

			if err != nil {
				lrw.statusCode = http.StatusInternalServerError
			}

			// Get the most up-to-date logger in case any other middlewares/security providers added to its context.
			l := zerolog.Ctx(r.Context())

			requestLogger := l.
				With().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Dur("elapsed_ms", time.Since(start)).
				Int("status_code", lrw.statusCode).
				Logger()

			if err != nil {
				requestLogger.Error().Interface("error", err).Msg("request handler panicked")
				panic(err)
			}

			requestLogger.Info().Msg("request handled")
		}()

		next.ServeHTTP(lrw, r)
	})
}
