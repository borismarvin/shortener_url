package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	api "github.com/borismarvin/shortener_url.git/internal/app/storage/api/errors"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/api/model"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/local"
	"go.uber.org/zap"
)

type FileStorage struct {
	model.Storage

	mutex sync.RWMutex

	file     *os.File
	fileName string
	encoder  *json.Encoder

	cache  local.LocalStorage
	lastID uint

	IsTemp bool
}

// Creates a new concurrent map
func NewFileStorage(fileName string) (*FileStorage, error) {
	if fileName == "" {
		zap.L().Info("storage was created successfully without keeping URL on disk")
		return &FileStorage{
			cache:  *local.NewLocalStorage(0),
			lastID: 0,
		}, nil
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage := &FileStorage{
		mutex:    sync.RWMutex{},
		file:     file,
		fileName: fileName,
		encoder:  json.NewEncoder(file),
		cache:    *local.NewLocalStorage(0),
		lastID:   0,
	}

	err = storage.fillCacheFromFile()
	if err != nil {
		return nil, err
	}

	zap.L().Info("storage was created successfully")

	return storage, nil
}

// Returns an element from the map
func (s *FileStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	s.mutex.RLock()
	res, ok := s.cache.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return entity.ErrorURLResponse(api.ErrShortURLNotFound)
	}

	return entity.OKURLResponse(res)
}

// Adds the given value under the specified key
//
// Returns `true` if element has been added to the storage.
func (s *FileStorage) AddURL(ctx context.Context, key, value entity.URL) entity.Response {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.cache.Get(key); ok {
		return entity.ErrorResponse(api.ErrURLAlreadyExists)
	}

	if s.file == nil {
		s.cache.Add(key, value)
		s.lastID++
		return entity.OKResponse()
	}

	storageRec := &entity.URLRecord{
		ID:          s.lastID + 1,
		ShortURL:    key.Path,
		OriginalURL: value.String(),
	}

	err := s.encoder.Encode(&storageRec)
	if err != nil {
		outputErr := fmt.Errorf("error while encoding entity for file commit: %w", err)
		return entity.ErrorResponse(outputErr)
	}

	s.file.Sync()

	s.cache.Add(key, value)
	s.lastID = storageRec.ID

	return entity.OKResponse()
}

func (s *FileStorage) Close() entity.Response {
	s.file.Name()
	if strings.Contains(s.fileName, os.TempDir()) {
		err := os.Remove(s.fileName)
		if err != nil {
			return entity.ErrorResponse(err)
		}
	}

	return entity.OKResponse()
}

func (s *FileStorage) PingServer(ctx context.Context) entity.Response {
	if s.file == nil {
		return entity.ErrorResponse(api.ErrFileStorageNotOpen)
	}

	return entity.OKResponse()
}

// Fills cache from the DB storage file
func (s *FileStorage) fillCacheFromFile() error {
	s.file.Seek(0, 0)
	scanner := bufio.NewScanner(s.file)
	record := entity.URLRecord{}

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}

		key, err := entity.NewURL(record.ShortURL)
		if err != nil {
			return err
		}

		value, err := entity.NewURL(record.OriginalURL)
		if err != nil {
			return err
		}

		s.cache.Add(*key, *value)
		s.lastID = record.ID
	}

	return nil
}
