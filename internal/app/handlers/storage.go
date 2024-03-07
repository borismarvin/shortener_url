package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type FileStorage struct {
	mx            sync.Mutex
	storageReader *reader
	storageWriter *writer
}

type reader struct {
	file    *os.File
	decoder *json.Decoder
}

type writer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func NewFileStorage(filename string) (fs *FileStorage, err error) {
	fs = &FileStorage{}
	fs.storageReader, err = NewReader(filename)
	if err != nil {
		return nil, err
	}
	fs.storageWriter, err = NewWriter(filename)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (p *writer) Write(event *Item) error {
	return p.encoder.Encode(&event)
}

func (p *writer) Close() error {
	return p.file.Close()
}

func (c *reader) Read() (*Item, error) {
	event := &Item{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *reader) Close() error {
	return c.file.Close()
}

// Save - сохраняет ID и ссылку в файле
func (f *FileStorage) Save(hash string, url string) error {
	f.mx.Lock()
	defer f.mx.Unlock()

	a := Item{Hash: hash, URL: url}
	err := f.storageWriter.Write(&a)
	if err != nil {
		return err
	}

	return nil
}

// Find ищет в файле ссылку
func (f *FileStorage) Find(hash string) (link string, err error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	_, err = f.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("Ошибка при установке указателя в бд - %s\n", err)
		return "", err
	}

	for {
		item, err := f.storageReader.Read()

		if err != nil {
			fmt.Printf("Ошибка при чтении из бд - %s\n", err)
			return "", err
		}

		if item.Hash == hash {
			return item.URL, nil
		}
	}
}

// Item - структура для хранения ссылки в файле
type Item struct {
	Hash string
	URL  string
}
