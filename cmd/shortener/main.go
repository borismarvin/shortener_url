// iter5
package main

import (
	"flag"
	"net/http"
	"os"

	middlewares "github.com/borismarvin/shortener_url.git/internal/app/middlewares"

	"github.com/borismarvin/shortener_url.git/cmd/shortener/config"
	handlers "github.com/borismarvin/shortener_url.git/internal/app/handlers"
	"github.com/borismarvin/shortener_url.git/internal/app/logger"
	"github.com/go-chi/chi/v5"
)

func InitializeConfig(startAddr string, baseAddr string, dbPath string) config.Args {
	envStartAddr := os.Getenv("SERVER_ADDRESS")
	envBaseAddr := os.Getenv("BASE_ADDRESS")
	envDBPath := os.Getenv("FILE_STORAGE_PATH")

	flag.StringVar(&startAddr, "a", "localhost:8080", "HTTP server start address")
	flag.StringVar(&baseAddr, "b", "http://localhost:8080", "Base address")
	flag.StringVar(&dbPath, "f", "./db", "Database path")
	flag.Parse()

	if envStartAddr != "" {
		startAddr = envStartAddr
	}
	if envBaseAddr != "" {
		baseAddr = envBaseAddr
	}
	if envDBPath != "" {
		dbPath = envDBPath
	}
	flag.Parse()

	builder := config.NewGetArgsBuilder()
	args := builder.
		SetStart(startAddr).
		SetBase(baseAddr).
		SetDB(dbPath).Build()
	return *args
}

func main() {

	var startAddr, baseAddr, dbPath string
	r := router()
	args := InitializeConfig(startAddr, baseAddr, dbPath)

	logger.Initialize()
	handlers.BaseURL = args.BaseAddr
	handlers.Storage, _ = handlers.NewFileStorage(args.DBPath)

	http.ListenAndServe(args.StartAddr, r)
}

func router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middlewares.Decompress)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Post("/api/shorten", handlers.APICreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)

	return r
}
