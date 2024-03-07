package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/borismarvin/shortener_url.git/internal/app/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostUrl(t *testing.T) {
	type want struct {
		response   string
		statusCode int
	}

	tests := []struct {
		name   string
		url    string
		method string
		body   io.Reader
		want   want
	}{
		{
			name:   "Получение короткой ссылки по полной",
			url:    "/",
			method: http.MethodPost,
			body:   strings.NewReader("http://ya.ru?x=fljdlfsdf&y=rweurowieur&z=sdkfhsdfisdf"),
			want: want{
				statusCode: http.StatusCreated,
				response:   "http://127.0.0.1:8080/d41d8cd98f00b204e9800998ecf8427e",
			},
		},
		{
			name:   "Получение полной ссылки по короткой",
			url:    "/d41d8cd98f00b204e9800998ecf8427e",
			method: http.MethodGet,
			body:   nil,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
			},
		},
	}

	handlers.Storage, _ = handlers.NewFileStorage("./db_test")
	handlers.BaseURL = "http://127.0.0.1:8080"

	r := router()
	ts := httptest.NewServer(r)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf(tt.name)

			response, body := testRequest(t, ts, tt.method, tt.url)

			defer response.Body.Close()

			assert.Equal(t, tt.want.statusCode, response.StatusCode)

			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, body)
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
