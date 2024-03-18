package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	shortenerErrors "github.com/borismarvin/shortener_url.git/internal/app/errors"
	"github.com/borismarvin/shortener_url.git/internal/app/middlewares"
	"github.com/borismarvin/shortener_url.git/internal/app/storage"
	"github.com/borismarvin/shortener_url.git/internal/app/types"
	"github.com/borismarvin/shortener_url.git/internal/app/utils"
	"github.com/go-chi/chi/v5"
)

var BaseURL string

// url для сокращения
type url struct {
	URL string `json:"url"`
}

// batchURL в пакетной обработке
type batchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// shortenBatchURL сокращенный урл в пакетной обработке
type shortenBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Сокращенный url
type response struct {
	URL string `json:"result"`
}

// URL пользователя
type userURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	originalURL, _ := ioutil.ReadAll(r.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("CreateShortURLHandler. %s", err)
		}
	}(r.Body)

	uuid := middlewares.UserSignedCookie.UUID
	hash, shortURL := utils.GetShortURL(string(originalURL))

	url := &types.Item{
		UUID:     uuid,
		Hash:     hash,
		URL:      string(originalURL),
		ShortURL: shortURL,
	}

	err := storage.Storage.Save(url)

	// Если такой url уже есть - отдаем соответствующий статус
	if errors.Is(err, shortenerErrors.ErrURLConflict) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(url.ShortURL))
		return
	}

	// Другие ошибки при сохранении в хранилище
	if err != nil {
		log.Printf("CreateShortURLHandler. Не удалось сохранить урл в хранилище. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortURL))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	exist, url, err := storage.Storage.FindByHash(hash)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url.URL))
}

// APICreateShortURLHandler Api для создания короткого урла
func APICreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	u := url{}

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := middlewares.UserSignedCookie.UUID
	hash, shortURL := utils.GetShortURL(string(u.URL))

	url := &types.Item{
		UUID:     uuid,
		Hash:     hash,
		URL:      u.URL,
		ShortURL: shortURL,
	}

	err := storage.Storage.Save(url)

	// Если такой url уже есть - отдаем соответствующий статус
	if errors.Is(err, shortenerErrors.ErrURLConflict) {
		resp, _ := json.Marshal(response{URL: url.ShortURL})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(resp)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	resp, _ := json.Marshal(response{URL: url.ShortURL})

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// APICreateShortURLBatchHandler Api для создания коротких урлов пачками
func APICreateShortURLBatchHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData []batchURL

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&incomingData); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var urls []*types.Item
	var resp []*shortenBatchURL
	uuid := middlewares.UserSignedCookie.UUID

	for _, url := range incomingData {
		shortURL := fmt.Sprintf("%s/%s", BaseURL, url.CorrelationID)

		urls = append(urls, &types.Item{
			UUID:     uuid,
			Hash:     url.CorrelationID,
			URL:      url.OriginalURL,
			ShortURL: shortURL,
		})
		resp = append(resp, &shortenBatchURL{
			CorrelationID: url.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	response, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// PingHandler проверяет соединение с базой
func PingHandler(w http.ResponseWriter, r *http.Request) {
	err := storage.Storage.Ping()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}
