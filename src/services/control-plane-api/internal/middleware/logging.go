package middleware

import (
	"context"
	"log/slog"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type loggerKey struct{}

func Logger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.With(
				"requestId", chimiddleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"route", r.Pattern)

			ctx := context.WithValue(r.Context(), loggerKey{}, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}
