package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: nil,
	}
}

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	contentTypes := c.w.Header().Values("Content-Type")
	if slices.Contains(contentTypes, "application/json") || slices.Contains(contentTypes, "text/html") {
		c.w.Header().Set("Content-Encoding", "gzip")
		c.zw = gzip.NewWriter(c.w)
		lenBuf, err := c.zw.Write(p)
		if err != nil {
			return 0, err
		}
		err = c.Close()
		if err != nil {
			return 0, err
		}
		return lenBuf, err
	} else {
		c.zw = nil
		return c.w.Write(p)
	}
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	if c.zw != nil {
		return c.zw.Close()
	}
	return nil
}

type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
func GzipMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			//allAcceptEncodingHeaders := strings.Split(r.Header.Values("Accept-Encoding")[0], ", ")
			var allAcceptEncodingSlice []string
			allAcceptEncodingHeaders := r.Header.Values("Accept-Encoding")
			if len(allAcceptEncodingHeaders) > 0 {
				allAcceptEncodingSlice = strings.Split(allAcceptEncodingHeaders[0], ", ")
			}
			if slices.Contains(allAcceptEncodingSlice, "gzip") {
				cw := NewCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncodings := r.Header.Values("Content-Encoding")
			if slices.Contains(contentEncodings, "gzip") {
				cr, err := NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}
			h.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
