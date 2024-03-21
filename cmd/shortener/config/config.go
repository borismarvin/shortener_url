package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	NetAddr           string `env:"SERVER_ADDRESS"`
	BaseURIPrefix     string `env:"BASE_URL"`
	LogLevel          string `env:"LOG_LEVEL"`
	DBFileStoragePath string `env:"FILE_STORAGE_PATH"`
	DBStorageConnect  string `env:"DATABASE_DSN"`
}

func InitConfig() (config Config) {
	flag.StringVar(&config.NetAddr, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&config.BaseURIPrefix, "b", "http://localhost:8080", "base output short URL")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.DBFileStoragePath, "f", "/tmp/short-url-db.json", "database storage path")
	flag.StringVar(&config.DBStorageConnect, "d", "", "database credentials in format: host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err.Error())
	}

	return
}
