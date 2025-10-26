package config

import "flag"

var flagRunBaseUrl string
var flagRunPrefixUrl string

type Config struct {
	BaseUrl, PrefixUrl string
}

func NewConfig() *Config {
	flag.StringVar(&flagRunBaseUrl, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagRunPrefixUrl, "b", "http://localhost:8080", "url prefix for links")
	flag.Parse()

	return &Config{
		BaseUrl:   flagRunBaseUrl,
		PrefixUrl: flagRunPrefixUrl,
	}
}
