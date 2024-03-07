package handlers

import (
	"bufio"
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
	currentID     int // Текущий индекс записи
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

	if ok, _ := f.IsEmpty(); ok {
		f.currentID = 0
	} else {
		lines, err := f.CountLines()
		if err != nil {
			return err
		}
		f.currentID = lines + 1
	}

	a := Item{UUID: f.currentID, ShortURL: hash, OriginalURL: url}
	err := f.storageWriter.Write(&a)
	if err != nil {
		return err
	}

	return nil
}

// Find searches for a URL in the file based on the given hash
func (f *FileStorage) Find(hash string) (link string, err error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	_, err = f.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("error seeking in the database: %s", err)
	}

	for {
		item, err := f.storageReader.Read()
		if err == io.EOF { // No more items in the file
			return "", fmt.Errorf("link not found for hash: %s", hash)
		} else if err != nil {
			return "", fmt.Errorf("error reading from the database: %s", err)
		}

		if item.ShortURL == hash {
			return item.OriginalURL, nil
		}
	}
}

// IsEmpty проверяет, пуст ли файл
func (f *FileStorage) IsEmpty() (bool, error) {
	fileInfo, err := f.storageReader.file.Stat()
	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func (f *FileStorage) CountLines() (int, error) {
	_, err := f.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("ошибка при поиске в базе данных: %s", err)
	}

	count := 0
	scanner := bufio.NewScanner(f.storageReader.file)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("ошибка при сканировании файла: %s", err)
	}

	return count, nil
}

// Item - структура для хранения ссылки в файле
type Item struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
