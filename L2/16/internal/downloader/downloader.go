package downloader

import (
	"fmt"
	"net/http"
	"time"
)

const (
	requestTimeout = 10 * time.Second
	userAgent      = "Mozilla/5.0 (compatible; Go-Site-Mirror/1.0;)"
)

// Fetch выполняет HTTP GET-запрос с таймаутом и возвращает тело ответа.
func Fetch(rawURL string) (*http.Response, error) {
	client := &http.Client{
		Timeout: requestTimeout,
	}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("неверный статус ответа: %s", resp.Status)
	}

	return resp, nil
}
