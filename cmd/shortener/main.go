// iter5
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/borismarvin/shortener_url.git/internal/app"
	middlewares "github.com/borismarvin/shortener_url.git/internal/app/middlewares"
	"github.com/caarlos0/env"

	"github.com/borismarvin/shortener_url.git/cmd/shortener/config"
	handlers "github.com/borismarvin/shortener_url.git/internal/app/handlers"
	"github.com/borismarvin/shortener_url.git/internal/app/logger"
	storage "github.com/borismarvin/shortener_url.git/internal/app/storage"
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

	r := router()

	logger.Initialize()

	// Переменные окружения в конфиг
	err := env.Parse(&app.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Параметры командной строки в конфиг
	flag.StringVar(&app.Cfg.ServerAddress, "a", app.Cfg.ServerAddress, "Адрес для запуска сервера")
	flag.StringVar(&app.Cfg.ServerPort, "server-port", app.Cfg.ServerPort, "Порт сервера")
	flag.StringVar(&app.Cfg.BaseURL, "b", app.Cfg.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&app.Cfg.DBPath, "f", app.Cfg.DBPath, "Путь к файлу с ссылками")
	flag.StringVar(&app.Cfg.DatabaseDsn, "d", app.Cfg.DatabaseDsn, "Строка с адресом подключения к БД")
	flag.Parse()

	log.Printf("Starting server on %s", app.Cfg.ServerAddress)
	log.Println(app.Cfg)

	// инициируем хранилище
	err = storage.New(&app.Cfg)
	if err != nil {
		log.Printf("Не удалось инициировать хранилище. %s", err)
		return
	}

	// запускаем сервер
	err = http.ListenAndServe(app.Cfg.ServerAddress, r)
	if err != nil {
		log.Printf("Не удалось запустить сервер. %s", err)
		return
	}
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
