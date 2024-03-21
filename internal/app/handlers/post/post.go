package post

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	storage "github.com/borismarvin/shortener_url.git/internal/app/storage/api/errors"
)

const (
	maxEncodedSize = 8
	timeout        = 3 * time.Second
)

type URLSaver interface {
	AddURL(ctx context.Context, key, value entity.URL) entity.Response
}

func postURLProcessing(saver URLSaver, ctx context.Context, inputURL, baseURIPrefix string) (string, error) {
	var shortURL *entity.URL

	userURL := entity.ParseURL(inputURL)
	added := false

	encodedURL := base64.StdEncoding.EncodeToString([]byte(inputURL))
	availableURLCount := len(encodedURL) / maxEncodedSize
	for i := 0; i < availableURLCount-1; i++ {
		shortURL = entity.ParseURL(encodedURL[(maxEncodedSize * i):(maxEncodedSize * (i + 1))])
		resp := saver.AddURL(ctx, *shortURL, *userURL)
		if resp.Status == entity.StatusOK {
			added = true
			break
		} else if !errors.Is(resp.Error, storage.ErrURLAlreadyExists) {
			return "", resp.Error
		}
	}

	if !added {
		return "", nil
	}

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL), nil
}
