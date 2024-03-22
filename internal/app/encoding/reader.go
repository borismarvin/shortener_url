package encoding

import (
	"compress/gzip"
	"io"
)

type compressReader struct {
	reader   io.ReadCloser
	gzReader *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		reader:   r,
		gzReader: gzR,
	}, nil
}

func (c *compressReader) Read(p []byte) (int, error) {
	return c.gzReader.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}

	return c.gzReader.Close()
}
