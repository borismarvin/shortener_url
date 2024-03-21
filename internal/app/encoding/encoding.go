package encoding

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const (
	gzipEncodingFormat = "gzip"
)

var (
	validContentTypes = map[string]struct{}{
		"application/json":         {},
		"text/html; charset=utf-8": {},
		"application/x-gzip":       {},
	}
)

func isEncodingContentType(content string) bool {
	if _, exists := validContentTypes[content]; exists {
		return true
	}

	return false
}

func GzipMiddleware(h http.Handler) http.Handler {
	gzipFn := func(writer http.ResponseWriter, req *http.Request) {
		ow := writer

		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, gzipEncodingFormat)
		if supportGzip {
			zap.L().Debug("sending gzip encoded message")
			cw := newCompressWriter(writer)
			cw.writer.Header().Set("Content-Encoding", "gzip")
			ow = cw
			defer cw.Close()
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		contentType := req.Header.Get("Content-Type")
		sendGzip := strings.Contains(contentEncoding, gzipEncodingFormat)
		if sendGzip && isEncodingContentType(contentType) {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				zap.L().Error("invalid compress reader creation", zap.Error(err))
				ow.WriteHeader(http.StatusInternalServerError)
				return
			}

			zap.L().Debug("obtained gzip encoded message")
			req.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, req)
	}

	return http.HandlerFunc(gzipFn)
}
