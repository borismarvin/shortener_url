// iter4
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/borismarvin/shortener_url.git/cmd/shortener/config"
	"github.com/gorilla/mux"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const keyLength = 6

type idToURLMap struct {
	links map[string]string
	id    string
	base  string
}

func main() {
	startAddr := flag.String("a", "localhost:8080", "HTTP server start address")
	baseAddr := flag.String("b", "http://localhost:8080", "Base address")
	flag.Parse()
	builder := config.NewGetArgsBuilder()
	args := builder.
		SetStart(*startAddr).
		SetBase(*baseAddr).Build()
	shortener := idToURLMap{
		links: make(map[string]string),
		base:  args.BaseAddr,
	}
	shortener.id = generateID()
	r := mux.NewRouter()
	shortenedURL := fmt.Sprintf("/%s", shortener.id)
	r.HandleFunc(shortenedURL, shortener.handleRedirect)
	r.HandleFunc("/", shortener.handleShortenURL)
	http.Handle("/", r)
	http.ListenAndServe(args.StartAddr, r)
}

func (iu idToURLMap) handleShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
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

func (iu idToURLMap) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
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
	return string(body), nil

}

func generateID() string {
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(shortKey)
}
