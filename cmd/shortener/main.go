// iter5
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/borismarvin/shortener_url.git/internal/app"
	middlewares "github.com/borismarvin/shortener_url.git/internal/app/middlewares"
	storage "github.com/borismarvin/shortener_url.git/internal/app/storage"
	"github.com/caarlos0/env"

	handlers "github.com/borismarvin/shortener_url.git/internal/app/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := Router()

	// Логер
	flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer flog.Close()

	//log.SetOutput(flog)

	// Переменные окружения в конфиг
	err = env.Parse(&app.Cfg)
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

func Router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.Decompress)

	r.Post("/", handlers.CreateShortURLHandler)
	r.Get("/ping", handlers.PingHandler)
	r.Post("/api/shorten", handlers.APICreateShortURLHandler)
	r.Get("/{hash}", handlers.GetShortURLHandler)

	return r
}
