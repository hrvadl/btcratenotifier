package cfg

import (
	"fmt"

	"github.com/caarlos0/env/v11"
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
	ExchangeServiceBaseURL string `env:"API_BASE_URL,required,notEmpty"`
	ExchangeServiceToken   string `env:"API_KEY,required,notEmpty"`
	LogLevel               string `env:"LOG_LEVEL,required,notEmpty"`
	Port                   string `env:"PORT,required,notEmpty"`
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
	var cfg Config
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "EXCHANGE_"}); err != nil {
		return nil, fmt.Errorf("%s failed: %w", operation, err)
	}
	return &cfg, nil
}
