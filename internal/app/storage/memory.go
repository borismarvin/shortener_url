package storage

import (
	"fmt"

	"github.com/borismarvin/shortener_url.git/internal/app/errors"
	"github.com/borismarvin/shortener_url.git/internal/app/types"
	"github.com/borismarvin/shortener_url.git/internal/app/utils"
)

type MemoryRepository struct {
	items map[string]*types.URL
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: map[string]*types.URL{},
	}
}

func (r *MemoryRepository) Save(url *types.URL) error {
	hash, _ := utils.GetShortURL(url.URL)

	// Дубли не храним
	if _, exist := r.items[hash]; !exist {
		r.items[hash] = url
		return nil
	} else {
		return fmt.Errorf("%w", errors.ErrURLConflict)
	}
}

func (r *MemoryRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	exist = false
	url = nil
	err = nil

	for _, item := range r.items {
		if item.Hash == hash {
			url = item
			exist = true
		}
	}

	return
}

func (r *MemoryRepository) FindByUUID(uuid string) (urls map[string]*types.URL, err error) {
	urls = map[string]*types.URL{}
	err = nil

	for _, item := range r.items {
		if item.UUID == uuid {
			urls[item.Hash] = item
		}
	}

	return
}
