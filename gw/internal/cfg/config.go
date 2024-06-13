package cfg

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const operation = "config parsing"

const (
	subServiceAddrEnvKey = "SUB_ADDR"
	rateWatchAddrEnvKey  = "RATE_WATCH_ADDR"
	logLevelEnvKey       = "GATEWAY_LOG_LEVEL"
	addrEnvKey           = "GATEWAY_ADDR"
)

// Config struct represents application config,
// which is used application-wide.
type Config struct {
	SubAddr         string `env:"SUB_ADDR,required,notEmpty"`
	RateWatcherAddr string `env:"RATE_WATCH_ADDR,required,notEmpty"`
	Addr            string `env:"GATEWAY_ADDR,required,notEmpty"`
	LogLevel        string `env:"GATEWAY_LOG_LEVEL,required,notEmpty"`
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
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("%s failed: %w", operation, err)
	}
	return &cfg, nil
}
