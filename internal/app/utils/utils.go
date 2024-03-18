package utils

import (
	"crypto/md5"
	"fmt"
)

var BaseURL string

// GetShortURL создает короткий урл из полного и возвращает хеш
func GetShortURL(value string) (hash string, shortURL string) {
	h := md5.New()
	h.Write([]byte(value))

	hash = fmt.Sprintf("%x", h.Sum(nil))
	shortURL = fmt.Sprintf("%s/%x", BaseURL, h.Sum(nil))

	return
}
