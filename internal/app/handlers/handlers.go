package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/borismarvin/shortener_url.git/internal/app"
	"github.com/borismarvin/shortener_url.git/internal/app/middlewares"
	storage "github.com/borismarvin/shortener_url.git/internal/app/storage"
	"github.com/borismarvin/shortener_url.git/internal/app/types"
	"github.com/borismarvin/shortener_url.git/internal/app/utils"
	"github.com/go-chi/chi/v5"
)

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
	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	uuid := middlewares.UserSignedCookie.UUID
	hash, shortURL := utils.GetShortURL(string(originalURL))

	url := &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      string(originalURL),
		ShortURL: shortURL,
	}

	err = storage.Storage.Save(url)

	// Если такой url уже есть - отдаем соответствующий статус
	if err != nil {
		fmt.Println("Ошибка при чтении тела запроса")
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortURL))
}

// GetShortURLHandler — возвращает полный урл по короткому.
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	exist, url, err := storage.Storage.FindByHash(hash)

	if !exist {
		fmt.Printf("Невозможно найти сслыку по хэшу - %s: %s", hash, err)
	}

	if err != nil {
		fmt.Printf("Ошибка поиска по хэшу - %s: %s", hash, err)
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

	url := &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      u.URL,
		ShortURL: shortURL,
	}

	err := storage.Storage.Save(url)

	if err != nil {
		fmt.Printf("Ошибка сохранения url - %s:", err)
	}

	resp, _ := json.Marshal(response{URL: url.ShortURL})

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// PingHandler проверяет соединение с базой
func PingHandler(w http.ResponseWriter, r *http.Request) {
	err := storage.Storage.Ping()

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

// ShortenMultipleUrl — принимающий в теле запроса множество URL для сокращения в формате:
//
//	{
//		"correlation_id": "<строковый идентификатор>",
//		"original_url": "<URL для сокращения>"
//	},
func ShortenMultipleURL(w http.ResponseWriter, r *http.Request) {
	var resp []*shortenBatchURL
	var urls []*types.URL
	var urlArray []batchURL
	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&urlArray); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := middlewares.UserSignedCookie.UUID
	for _, u := range urlArray {
		shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, u.CorrelationID)

		urls = append(urls, &types.URL{
			UUID:     uuid,
			Hash:     u.CorrelationID,
			URL:      u.OriginalURL,
			ShortURL: shortURL,
		})
		resp = append(resp, &shortenBatchURL{
			CorrelationID: u.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	err := storage.Storage.SaveBatch(urls)
	if err != nil {
		fmt.Println(err)
	}

	response, _ := json.Marshal(resp)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
