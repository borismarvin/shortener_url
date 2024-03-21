package post

import (
	"context"
	"io"
	"net/http"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	"github.com/borismarvin/shortener_url.git/internal/app/handlers/errors"
	"go.uber.org/zap"
)

// Processes POST request by URL within http://localhost:8080/id URL format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandlerURL(saver URLSaver, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		zap.L().Debug("POST handler URL processing")

		inputURL, err := io.ReadAll(req.Body)
		defer req.Body.Close()

		if err != nil {
			zap.L().Error(errors.CannotProcessURL, zap.Error(err))
			http.Error(writer, errors.WrongURLFormat, http.StatusBadRequest)
			return
		}

		if ok := entity.IsValidURL(string(inputURL)); !ok {
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

		outputURL, err := postURLProcessing(saver, ctx, string(inputURL), baseURIPrefix)
		if err != nil || outputURL == "" {
			zap.L().Error("could not create a short URL", zap.String("error", err.Error()))
			http.Error(writer, errors.InternalServerError, http.StatusInternalServerError)
			return
		}

		zap.L().Info("url has been created succeessfully", zap.String("output url", outputURL))

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusCreated)
		io.WriteString(writer, outputURL)
	}
}
