package get

import (
	"context"
	"net/http"
	"time"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(ctx context.Context, key entity.URL) entity.URLResponse
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func GetURLHandler(getter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		shortURL := chi.URLParam(req, "url")

		ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
		defer cancel()

		resp := getter.GetURL(ctx, *entity.ParseURL(shortURL))
		if resp.Status == entity.StatusError {
			zap.L().Error(
				"error while getting url",
				zap.String("error", resp.Error.Error()),
				zap.String("short_url", shortURL),
				zap.String("decoded_url", resp.URL.String()),
			)

			http.Error(writer, errors.ShortURLNotInDB, http.StatusBadRequest)
			return
		}

		zap.L().Info("url has been decoded succeessfully", zap.String("decoded url", resp.URL.String()))

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.Header().Set("Location", resp.URL.String())
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}
