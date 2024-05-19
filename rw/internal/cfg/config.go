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
	portEnvKey                   = "EXCHANGE_PORT"
)

// Config struct represents application config,
// which is used application-wide.
type Config struct {
	ExchangeServiceBaseURL string
	ExchangeServiceToken   string
	LogLevel               string
	Port                   string
}

// Must is a handly wrapper around return results from
// the NewFromEnv() function, which will panic in case of error.
// Should be called only in main function, when we don't need
// to handle errors.
func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

// NewFromEnv parses the environment variables into
// the Config struct. Returns an error if any of required variables
// is missing or contains invalid value.
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

	port := os.Getenv(portEnvKey)
	if port == "" {
		return nil, fmt.Errorf("%s: port can't be empty", operation)
	}

	return &Config{
		ExchangeServiceBaseURL: apiURL,
		ExchangeServiceToken:   apiKey,
		LogLevel:               logLevel,
		Port:                   port,
	}, nil
}
