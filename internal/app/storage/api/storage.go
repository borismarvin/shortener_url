package api

import (
	"github.com/borismarvin/shortener_url.git/internal/app/config"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/api/model"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/file"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/local"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/postgres"
	"go.uber.org/zap"
)

func InitStorage(config config.Config) (model.Storage, error) {
	var db model.Storage
	var err error

	if len(config.DBStorageConnect) > 0 {
		zap.L().Info("init postgres storage")
		db, err = postgres.NewPostgresStorage(config.DBStorageConnect)
	} else if len(config.DBFileStoragePath) > 0 {
		zap.L().Info("init file storage")
		db, err = file.NewFileStorage(config.DBFileStoragePath)
	} else {
		zap.L().Info("init local storage")
		db = local.NewTSLocalStorage(0)
		err = nil
	}

	return db, err
}
