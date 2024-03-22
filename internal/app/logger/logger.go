package logger

import (
	"net/http"
	"time"

	"github.com/borismarvin/shortener_url.git/internal/app/config"
	"go.uber.org/zap"
)

func Initialize(config config.Config) error {
	lvl, err := zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(zl)

	return nil
}

func LoggerMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{
			statusCode: 0,
			size:       0,
		}
		writer := logResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}
		h.ServeHTTP(&writer, r)

		duration := time.Since(start)

		zap.L().Info(
			"got incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", respData.statusCode),
			zap.Int("size", respData.size),
		)
	}

	return http.HandlerFunc(logFn)
}
