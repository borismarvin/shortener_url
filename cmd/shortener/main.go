// iter1
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const keyLength = 6

type idToURLMap struct {
	links map[string]string
	id    string
}

func main() {
	shortener := idToURLMap{
		links: make(map[string]string),
	}
	shortener.id = generateID()

	shortenedURL := fmt.Sprintf("/%s", shortener.id)
	http.HandleFunc(shortenedURL, shortener.handleRedirect)
	http.HandleFunc("/", shortener.handleShortenURL)
	http.ListenAndServe(":8080", nil)
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

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", id)

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
