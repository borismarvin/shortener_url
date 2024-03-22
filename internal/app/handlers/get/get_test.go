package get

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/borismarvin/shortener_url.git/internal/app/config"
	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/errors"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/get/mock"
	"github.com/borismarvin/shortener_url.git/internal/app/logger"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLGetter(ctrl)

	type want struct {
		statusCode  int
		contentType string
		location    string
		resp        entity.URLResponse
		message     string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "correct input data",
			request: "aHR0cHM6",

			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				contentType: "text/plain; charset=utf-8",
				location:    "https://practicum.yandex.ru/",
				resp:        makeOKURLResponse("https://practicum.yandex.ru/"),
				message:     "",
			},
		},
		{
			name:    "request without id",
			request: "",

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
				resp:        entity.ErrorURLResponse(fmt.Errorf("")),
				message:     errors.ShortURLNotInDB + "\n",
			},
		},
		{
			name:    "missing URL",
			request: "/fsdfuytu",

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
				resp:        entity.ErrorURLResponse(fmt.Errorf("")),
				message:     errors.ShortURLNotInDB + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/{url}", nil)
			writer := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("url", test.request)

			s.EXPECT().GetURL(gomock.Any(), gomock.Any()).
				Return(test.want.resp)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			handler := GetURLHandler(s)
			handler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.message, string(userResult))
		})
	}
}

func TestGetPingDBHandler(t *testing.T) {
	cnf := config.InitConfig()

	logger.Initialize(cnf)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockStoragePinger(ctrl)

	type want struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name string
		resp entity.Response
		want want
	}{
		{
			name: "successfull ping",
			resp: entity.OKResponse(),

			want: want{
				statusCode: http.StatusOK,
				err:        nil,
			},
		},
		{
			name: "fallen ping",
			resp: entity.ErrorResponse(fmt.Errorf("")),

			want: want{
				statusCode: http.StatusInternalServerError,
				err:        context.DeadlineExceeded,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/{url}", nil)
			writer := httptest.NewRecorder()

			s.EXPECT().
				PingServer(gomock.Any()).
				Return(test.resp)

			// pingDB(s, writer, request)
			handler := GetPingDB(s)
			handler(writer, request)

			res := writer.Result()

			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}

func makeOKURLResponse(URL string) entity.URLResponse {
	return entity.OKURLResponse(*entity.ParseURL(URL))
}
