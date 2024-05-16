package cfg

import (
	"fmt"
	"os"
)

const operation = "config parsing"

const (
	exchangeServiceBaseURLEnvKey = "EXCHANGE_API_BASE_URL"
	exchangeServiceTokenEnvKey   = "EXCHANGE_API_KEY"
	logLevelEnvKey               = "EXCHANGE_LOG_LEVEL"
)

type Config struct {
	ExchangeServiceBaseURL string
	ExchangeServiceToken   string
	LogLevel               string
}

func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewFromEnv() (*Config, error) {
	apiKey := os.Getenv(exchangeServiceTokenEnvKey)
	if apiKey == "" {
		return nil, fmt.Errorf("%s: api key for exchange service can't be empty", operation)
	}

	apiURL := os.Getenv(exchangeServiceBaseURLEnvKey)
	if apiURL == "" {
		return nil, fmt.Errorf("%s: api url for exchange service can't be empty", operation)
	}

	logLevel := os.Getenv(logLevelEnvKey)
	if logLevel == "" {
		return nil, fmt.Errorf("%s: log level can't be empty", operation)
	}

	return &Config{
		ExchangeServiceBaseURL: apiURL,
		ExchangeServiceToken:   apiKey,
		LogLevel:               logLevel,
	}, nil
}
