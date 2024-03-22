package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	"os"
	"strings"

	"github.com/borismarvin/shortener_url.git/internal/app/config"
	"github.com/borismarvin/shortener_url.git/internal/app/logger"
	"go.uber.org/zap"
)

func main() {
	cnf := config.InitConfig()

	err := logger.Initialize(cnf)
	if err != nil {
		panic(err.Error())
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var data = "https://practicum.yandex.ru/"
	url, err := postRequest(client, data, cnf.NetAddr)
	if err != nil {
		log.Println("Post request error: : %w", err)
		zap.L().Fatal("post request error", zap.Error(err))
		panic(err)
	}

	getRequest(client, url, cnf.BaseURIPrefix)
}

func readFromConsole() string {
	// приглашение в консоли
	fmt.Println("Input URL")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	data, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	data = strings.TrimSuffix(data, "\n")

	return data
}

func postRequest(client *http.Client, data, netAddr string) (string, error) {
	url := fmt.Sprintf("http://%s", netAddr)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(data)))
	if err != nil {
		zap.L().Error("request couldn't be created", zap.String("method", "POST"), zap.Error(err))
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		zap.L().Error("request failed", zap.String("method", "POST"), zap.Error(err))
		return "", err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		zap.L().Error("failed to read response body", zap.String("method", "POST"), zap.Error(err))
		return "", err
	}

	zap.L().Info(
		"request output",
		zap.String("method", "POST"),
		zap.String("body", string(bodyBytes)),
		zap.Int("status", response.StatusCode),
	)

	return path.Base(string(bodyBytes)), nil
}

func getRequest(client *http.Client, url, baseURIPrefix string) {
	requestURL := fmt.Sprintf("%s/%s", baseURIPrefix, url)
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		zap.L().Error("request couldn't be created", zap.String("method", "GET"), zap.Error(err))
		return
	}

	response, err := client.Do(request)
	if err != nil {
		zap.L().Error("request failed", zap.String("method", "GET"), zap.Error(err))
		return
	}
	defer response.Body.Close()

	zap.L().Info(
		"request output",
		zap.String("method", "GET"),
		zap.String("location", response.Header.Get("Location")),
		zap.Int("status", response.StatusCode),
	)
}
