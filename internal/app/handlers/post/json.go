package post

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/errors"
	"github.com/borismarvin/shortener_url.git/internal/app/models"
	"go.uber.org/zap"
)

// Processes POST request by JSON within http://localhost:8080/api/shorten URL format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandlerJSON(saver URLSaver, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("POST handler JSON processing")

		inputRequest := &models.Request{}
		err := json.NewDecoder(req.Body).Decode(&inputRequest)
		defer req.Body.Close()
		if err != nil {
			zap.L().Error(errors.CannotProcessJSON, zap.Error(err))
			http.Error(writer, errors.WrongURLFormat, http.StatusBadRequest)
			return
		}

		if ok := entity.IsValidURL(inputRequest.URL); !ok {
			zap.L().Error(errors.WrongURLFormat, zap.Error(err))
			http.Error(writer, errors.WrongURLFormat, http.StatusBadRequest)
			return
		}

		if baseURIPrefix == "" {
			zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
			http.Error(writer, errors.InternalServerError, http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		outputURL, err := postURLProcessing(saver, ctx, inputRequest.URL, baseURIPrefix)
		if err != nil || outputURL == "" {
			zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
			http.Error(writer, errors.InternalServerError, http.StatusInternalServerError)
			return
		}

		response := &models.Response{
			URL: outputURL,
		}
		zap.L().Info("url has been created succeessfully", zap.String("output url", response.URL))

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(writer).Encode(response); err != nil {
			zap.L().Error("invalid response", zap.Any("response", response))
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			return
		}

		zap.L().Debug("sending HTTP 200 response")
	}
}
