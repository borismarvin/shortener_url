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

func InitializeConfig(startAddr string, baseAddr string, filePath string, database string) config.Args {
	envStartAddr := os.Getenv("SERVER_ADDRESS")
	envBaseAddr := os.Getenv("BASE_ADDRESS")
	envFilePath := os.Getenv("FILE_STORAGE_PATH")
	envDatabase := os.Getenv("DATABASE_DSN")

	flag.StringVar(&startAddr, "a", "localhost:8080", "HTTP server start address")
	flag.StringVar(&baseAddr, "b", "http://localhost:8080", "Base address")
	flag.StringVar(&filePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&database, "d", "host=localhost port=5432 user=postgres password=123 dbname=db sslmode=disable", "Database file")
	flag.Parse()

	if envStartAddr != "" {
		startAddr = envStartAddr
	}
	if envBaseAddr != "" {
		baseAddr = envBaseAddr
	}
	if envFilePath != "" {
		filePath = envFilePath
	}
	if envDatabase != "" {
		database = envDatabase
	}
	flag.Parse()

	builder := config.NewGetArgsBuilder()
	args := builder.
		SetStart(startAddr).
		SetBase(baseAddr).
		SetFile(filePath).
		SetDB(database).Build()
	return *args
}

func main() {

	var startAddr, baseAddr, filePath, dbPath string
	r := router()
	args := InitializeConfig(startAddr, baseAddr, filePath, dbPath)

	logger.Initialize()
	handlers.BaseURL = args.BaseAddr
	handlers.Storage, _ = handlers.NewFileStorage(args.FilePath)
	handlers.DSN = args.Database
	http.ListenAndServe(args.StartAddr, r)
}

func router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middlewares.Decompress)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Post("/api/shorten", handlers.APICreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)
	r.Get("/ping", handlers.PingHandler)

	return r
}
