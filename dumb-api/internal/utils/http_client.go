package utils

import (
	"net/http"
	"time"
)

// Create a shared HTTP client with timeouts, etc.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
