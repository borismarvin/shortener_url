package storage

import (
	"fmt"

	"github.com/borismarvin/shortener_url.git/internal/app/types"
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

	// Дубли не храним
	if _, exist := r.items[url.Hash]; !exist {
		r.items[url.Hash] = url
		return nil
	} else {
		return fmt.Errorf("url уже существует")
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
