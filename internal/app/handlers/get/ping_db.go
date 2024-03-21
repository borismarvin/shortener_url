package get

import (
	"context"
	"errors"
	"net/http"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	"go.uber.org/zap"
)

type StoragePinger interface {
	PingServer(ctx context.Context) entity.Response
}

func GetPingDB(pinger StoragePinger) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		response := pinger.PingServer(ctx)
		if response.Status == entity.StatusError {
			switch {
			case errors.Is(response.Error, context.Canceled):
				zap.L().Error("context canceled", zap.String("error", response.Error.Error()))
			case errors.Is(response.Error, context.DeadlineExceeded):
				zap.L().Error("context deadline exceeded", zap.String("error", response.Error.Error()))
			default:
				zap.L().Error(response.Error.Error())
			}
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		zap.L().Info("storage works after ping")

		writer.WriteHeader(http.StatusOK)
	}
}
