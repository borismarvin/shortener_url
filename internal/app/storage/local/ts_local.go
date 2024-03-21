package local

import (
	"context"
	"sync"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	api "github.com/borismarvin/shortener_url.git/internal/app/storage/api/errors"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/api/model"
)

type TSLocalStorage struct {
	model.Storage

	mutex sync.RWMutex
	urls  LocalStorage
}

func NewTSLocalStorage(size int) *TSLocalStorage {
	return &TSLocalStorage{
		urls: *NewLocalStorage(size),
	}
}

// Returns an element from the map
func (s *TSLocalStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	s.mutex.RLock()
	res, ok := s.urls.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return entity.ErrorURLResponse(api.ErrShortURLNotFound)
	}

	return entity.OKURLResponse(res)
}

// Adds the given value under the specified key
func (s *TSLocalStorage) AddURL(ctx context.Context, key, value entity.URL) entity.Response {
	s.mutex.Lock()
	s.urls.Add(key, value)
	s.mutex.Unlock()

	return entity.OKResponse()
}

func PingServer(ctx context.Context) entity.Response {
	return entity.OKResponse()
}

func Close() entity.Response {
	return entity.OKResponse()
}
