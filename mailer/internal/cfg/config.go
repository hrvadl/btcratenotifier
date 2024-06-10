package cfg

import (
	"fmt"
	"os"
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
	MailerToken string
	LogLevel    string
	Port        string
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
	logLevel := os.Getenv(logLevelEnvKey)
	if logLevel == "" {
		return nil, fmt.Errorf("%s: log level can't be empty", operation)
	}

	port := os.Getenv(portEnvKey)
	if port == "" {
		return nil, fmt.Errorf("%s: port can't be empty", operation)
	}

	mailerToken := os.Getenv(mailerTokenEnvKey)
	if mailerToken == "" {
		return nil, fmt.Errorf("%s: token can't be empty", operation)
	}

	return &Config{
		LogLevel:    logLevel,
		Port:        port,
		MailerToken: mailerToken,
	}, nil
}
