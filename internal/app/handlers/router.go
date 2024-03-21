package handlers

import (
	"github.com/borismarvin/shortener_url.git/internal/app/config"
	"github.com/borismarvin/shortener_url.git/internal/app/encoding"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/get"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/post"
	"github.com/borismarvin/shortener_url.git/internal/app/logger"
	storage "github.com/borismarvin/shortener_url.git/internal/app/storage/api/model"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(config config.Config, db storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.LoggerMiddleware)
	r.Use(encoding.GzipMiddleware)

	r.Post("/", post.PostHandlerURL(db, config.BaseURIPrefix))
	r.Post("/api/shorten", post.PostHandlerJSON(db, config.BaseURIPrefix))

	r.Get("/{url}", get.GetURLHandler(db))
	r.Get("/ping", get.GetPingDB(db))

	return r
}
