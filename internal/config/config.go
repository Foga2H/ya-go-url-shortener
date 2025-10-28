// internal/config/config.go
package config

import (
	"flag"
	"os"
)

var (
	flagRunBaseURL   = flag.String("a", "localhost:8080", "address and port to run server")
	flagRunPrefixURL = flag.String("b", "http://localhost:8080", "url prefix for links")
)

type Config struct {
	BaseURL   string
	PrefixURL string
}

func NewConfig() *Config {
	if envRunBaseUrl := os.Getenv("SERVER_ADDRESS"); envRunBaseUrl != "" {
		*flagRunBaseURL = envRunBaseUrl
	}

	if envRunPrefixURL := os.Getenv("BASE_URL"); envRunPrefixURL != "" {
		*flagRunPrefixURL = envRunPrefixURL
	}

	return &Config{
		BaseURL:   *flagRunBaseURL,
		PrefixURL: *flagRunPrefixURL,
	}
}

func NewConfigFrom(baseURL, prefixURL string) *Config {
	return &Config{
		BaseURL:   baseURL,
		PrefixURL: prefixURL,
	}
}
