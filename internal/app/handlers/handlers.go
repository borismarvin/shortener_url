package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/borismarvin/shortener_url.git/internal/app/storage"
	"github.com/borismarvin/shortener_url.git/internal/app/types"
	"github.com/borismarvin/shortener_url.git/internal/app/utils"

	"github.com/go-chi/chi/v5"
)

// url для сокращения
type url struct {
	URL string `json:"url"`
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length int = 10

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateUUID() string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// CreateShortURLHandler — создает короткий урл.
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hash, shortURL := utils.GetShortURL(string(originalURL))
	uuid := GenerateUUID()
	url := &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      string(originalURL),
		ShortURL: shortURL,
	}

	err = storage.Storage.Save(url)

	if err != nil {
		fmt.Println("url уже существует")
	}

	// Другие ошибки при сохранении в хранилище
	if err != nil {
		log.Printf("CreateShortURLHandler. Не удалось сохранить урл в хранилище. %s", err)
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
		fmt.Printf("Невозможно найти сслыку по хэшу - %s: %s", hash, err)
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

	hash, shortURL := utils.GetShortURL(string(u.URL))
	uuid := GenerateUUID()
	url := &types.URL{
		UUID:     uuid,
		Hash:     hash,
		URL:      u.URL,
		ShortURL: shortURL,
	}

	err := storage.Storage.Save(url)

	if err != nil {
		fmt.Println("url уже существует")
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

// првоеряет соединение с базой данных
func PingPong(w http.ResponseWriter, r *http.Request) {
	if CheckDBConn(DatabaseName) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}
