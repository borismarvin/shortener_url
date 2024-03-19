package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/borismarvin/shortener_url.git/internal/app/types"
)

func newWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func newReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func NewFileRepository(filename string) (r *FileRepository, err error) {
	r = &FileRepository{}
	r.storageReader, err = newReader(filename)
	if err != nil {
		return nil, err
	}
	r.storageWriter, err = newWriter(filename)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type writer struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *writer) Write(url *types.URL) error {
	return p.encoder.Encode(&url)
}

func (p *writer) Close() error {
	return p.file.Close()
}

type reader struct {
	file    *os.File
	decoder *json.Decoder
}

func (c *reader) Read() (*types.URL, error) {
	item := &types.URL{}
	if err := c.decoder.Decode(&item); err != nil {
		return nil, err
	}
	return item, nil
}

func (c *reader) Close() error {
	return c.file.Close()
}

type FileRepository struct {
	mx            sync.Mutex
	storageReader *reader
	storageWriter *writer
}

func (r *FileRepository) Save(url *types.URL) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	err := r.storageWriter.Write(url)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	_, err = r.storageReader.file.Seek(0, 0)
	if err != nil {
		return false, &types.URL{}, err
	}

	for {
		item, err := r.storageReader.Read()

		if err != nil {
			return false, nil, err
		}

		if item.Hash == hash {
			return true, item, nil
		}
	}
}

func (r *FileRepository) FindByUUID(uuid string) (urls map[string]*types.URL, err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	urls = map[string]*types.URL{}

	_, err = r.storageReader.file.Seek(0, 0)
	if err != nil {
		return map[string]*types.URL{}, err
	}

	for {
		item, err := r.storageReader.Read()

		if err != nil {
			break
		}

		if item.UUID == uuid {
			urls[item.Hash] = item
		}
	}

	return urls, nil
}
