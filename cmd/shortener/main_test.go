package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRedirect(t *testing.T) {
	m := idToURLMap{
		links: map[string]string{
			"123": "https://practicum.yandex.ru/",
		},
		id: "123",
	}
	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", m.id)
	req, err := http.NewRequest("GET", shortenedURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(m.handleRedirect)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	expectedLocation := m.links[m.id]
	if location := rr.Header().Get("Location"); location != expectedLocation {
		t.Errorf("handler returned unexpected location header: got %v want %v",
			location, expectedLocation)
	}
}

func TestHandleShortenURL(t *testing.T) {
	m := idToURLMap{
		links: map[string]string{
			"123": "https://practicum.yandex.ru/",
		},
		id: "123",
	}
	originalURL := m.links[m.id]
	body := strings.NewReader("https://practicum.yandex.ru/")
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	m.handleShortenURL(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusCreated)
	}

	expectedContentType := "text/plain"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected Content-Type header: got %v want %v",
			contentType, expectedContentType)
	}

	expectedURL := fmt.Sprintf("http://localhost:8080/%s", m.id)
	bodyBytes := rr.Body.Bytes()
	if string(bodyBytes) != expectedURL {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(bodyBytes), expectedURL)
	}

	if url := m.links[m.id]; url != originalURL {
		t.Errorf("handler failed to add URL to map: got %v want %v",
			url, originalURL)
	}
}
