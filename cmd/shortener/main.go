// iter5
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/borismarvin/shortener_url.git/cmd/shortener/config"
	"github.com/borismarvin/shortener_url.git/internal/app/logger"
	"github.com/gorilla/mux"
)

const (
	Unsupported ContentType = iota
	PlainText
	URLEncoded
	JSON
)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const keyLength = 6

type idToURLMap struct {
	links map[string]string
	id    string
	base  string
}

type SendData struct {
	URL string `json:"result"`
}
type GetData struct {
	Result string `json:"result"`
}

type ContentType int

var supportedTypes = []ContentTypes{
	{
		name: "text/plain",
		code: PlainText,
	},
	{
		name: "application/x-www-form-urlencoded",
		code: URLEncoded,
	},
	{
		name: "application/json",
		code: JSON,
	},
}

type ContentTypes struct {
	name string
	code ContentType
}

func InitializeConfig(startAddr string, baseAddr string) config.Args {
	envStartAddr := os.Getenv("SERVER_ADDRESS")
	envBaseAddr := os.Getenv("BASE_ADDRESS")

	flag.StringVar(&startAddr, "a", "localhost:8080", "HTTP server start address")
	flag.StringVar(&baseAddr, "b", "http://localhost:8080", "Base address")
	flag.Parse()

	if envStartAddr != "" {
		startAddr = envStartAddr
	}
	if envBaseAddr != "" {
		baseAddr = envBaseAddr
	}
	flag.Parse()

	builder := config.NewGetArgsBuilder()
	args := builder.
		SetStart(startAddr).
		SetBase(baseAddr).Build()
	return *args
}

func main() {

	var startAddr, baseAddr string

	args := InitializeConfig(startAddr, baseAddr)

	r := mux.NewRouter()

	shortener := idToURLMap{
		links: make(map[string]string),
		base:  args.BaseAddr,
	}
	logger.Initialize()
	shortener.id = generateID()
	shortenedURL := fmt.Sprintf("/%s", shortener.id)
	r.HandleFunc(shortenedURL, shortener.handleRedirect)
	r.HandleFunc("/", shortener.handleShortenURL)
	originalURL := "https://practicum.yandex.ru"
	http.Post(`http://localhost:8080/api/shorten`, `application/json`,
		// ключи указаны в разных регистрах, но данные сконвертируются правильно
		bytes.NewBufferString(`{"url":`+`"`+originalURL+`"}`))
	r.HandleFunc("/api/shorten", shortener.handleShortenURLJSON)
	http.Handle("/", r)
	http.ListenAndServe(args.StartAddr, logger.WithLogging(r))
}

func (iu idToURLMap) handleShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	url, err := decodeRequestBody(w, r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id := iu.id
	iu.links[id] = url

	shortenedURL := iu.base + "/" + id
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortenedURL))
}

func getContentTypeCode(name string) ContentType {
	for _, t := range supportedTypes {
		if name == t.name {
			return t.code
		}
	}
	return Unsupported
}
func (iu idToURLMap) handleShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if getContentTypeCode(contentType) == JSON {
		var result GetData
		var url SendData
		if r.Method == http.MethodPost {

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err = json.Unmarshal(buf.Bytes(), &url); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id := iu.id
			iu.links[id] = url.URL

			shortenedURL := iu.base + "/" + id
			result.Result = shortenedURL

		}
		resp, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (iu idToURLMap) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := iu.id
	originalURL, ok := iu.links[id]
	if !ok {
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
func decodeRequestBody(w http.ResponseWriter, r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
	}
	return string(body), err

}

func generateID() string {
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(shortKey)
}
