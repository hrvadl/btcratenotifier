package cfg

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const operation = "config parsing"

const (
	logLevelEnvKey    = "MAILER_LOG_LEVEL"
	portEnvKey        = "MAILER_PORT"
	mailerTokenEnvKey = "MAILER_API_KEY" // #nosec G101
)

// Config struct represents application config,
// which is used application-wide.
type Config struct {
	MailerToken string `env:"API_KEY,required,notEmpty"`
	LogLevel    string `env:"LOG_LEVEL,required,notEmpty"`
	Port        string `env:"PORT,required,notEmpty"`
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
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "MAILER_"}); err != nil {
		return nil, fmt.Errorf("%s failed: %w", operation, err)
	}
	return &cfg, nil
}
