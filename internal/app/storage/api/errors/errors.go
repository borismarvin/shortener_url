package api

import "errors"

var (
	// ErrEmptyDBStorageCred = errors.New("error in db storage connection credentials are empty")
	ErrShortURLNotFound   = errors.New("short url is not found in storage")
	ErrURLAlreadyExists   = errors.New("short url already exists in storage")
	ErrFileStorageNotOpen = errors.New("file storage is not open")
)
