package client

import (
	"net/http"
	"sync"
	"time"
)

var (
	httpClient *http.Client
	once       sync.Once
)

func CreateHttpClient(timeout int) *http.Client {
	once.Do(func() {
		httpClient = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	})
	return httpClient
}
