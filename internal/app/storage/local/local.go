package local

import "github.com/borismarvin/shortener_url.git/internal/app/entity"

type LocalStorage struct {
	urls map[entity.URL]entity.URL
}

func NewLocalStorage(size int) *LocalStorage {
	return &LocalStorage{
		urls: make(map[entity.URL]entity.URL, size),
	}
}

// Returns an element from the map
func (s *LocalStorage) Get(key entity.URL) (entity.URL, bool) {
	res, ok := s.urls[key]

	return res, ok
}

// Adds the given value under the specified key
func (s *LocalStorage) Add(key, value entity.URL) {
	s.urls[key] = value
}
