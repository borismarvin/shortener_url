package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/borismarvin/shortener_url.git/internal/app"
	"github.com/borismarvin/shortener_url.git/internal/app/storage"
	"github.com/borismarvin/shortener_url.git/internal/app/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var S suite

type suite struct {
	Server *httptest.Server
}

func setup() {
	app.Cfg = types.Config{
		BaseURL:       "http://localhost:8080",
		ServerPort:    "8080",
		ServerAddress: "localhost:8080",
		DBPath:        "./db_test",
	}

	storage.New(&app.Cfg)

	S = suite{
		Server: httptest.NewServer(Router()),
	}
}

func TestPostUrl(t *testing.T) {
	setup()
	defer S.Server.Close()
	defer storage.Storage.Drop()

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
			name:   "Получение короткой ссылки",
			url:    "/",
			method: http.MethodPost,
			body:   strings.NewReader("http://jwlqct1udntv.com/xr0cz5fshffj/pimnbpv/otw2im3fudstqi1"),
			want: want{
				statusCode: http.StatusCreated,
				response:   "http://localhost:8080/580c5ab5ef6a4f27b3da9956ae192f4f",
			},
		},
		{
			name:   "Получение полной ссылки",
			url:    "/580c5ab5ef6a4f27b3da9956ae192f4f",
			method: http.MethodGet,
			body:   nil,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf(tt.name)

			response, body := testRequest(t, tt.method, tt.url, tt.body)
			defer response.Body.Close()

			assert.Equal(t, tt.want.statusCode, response.StatusCode)

			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, body)
			}
		})
	}
}

func testRequest(t *testing.T, method string, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, S.Server.URL+path, body)
	require.NoError(t, err)

	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
