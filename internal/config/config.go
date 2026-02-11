// internal/config/config.go
package config

import (
	"flag"
	"os"
)

var (
	flagRunBaseURL   = flag.String("a", "localhost:8080", "address and port to run server")
	flagRunPrefixURL = flag.String("b", "http://localhost:8080", "url prefix for links")
	fileStoragePath  = flag.String("f", "./storage.json", "path to file storage")
)

type Config struct {
	BaseURL         string
	PrefixURL       string
	FileStoragePath string
}

func NewConfig() *Config {
	if envRunBaseURL := os.Getenv("SERVER_ADDRESS"); envRunBaseURL != "" {
		*flagRunBaseURL = envRunBaseURL
	}

	if envRunPrefixURL := os.Getenv("BASE_URL"); envRunPrefixURL != "" {
		*flagRunPrefixURL = envRunPrefixURL
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		*fileStoragePath = envFileStoragePath
	}

	return &Config{
		BaseURL:         *flagRunBaseURL,
		PrefixURL:       *flagRunPrefixURL,
		FileStoragePath: *fileStoragePath,
	}
}

func NewConfigFrom(baseURL, prefixURL string) *Config {
	return &Config{
		BaseURL:   baseURL,
		PrefixURL: prefixURL,
	}
}
