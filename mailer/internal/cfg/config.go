package cfg

import (
	"fmt"
	"os"
)

const operation = "config parsing"

const (
	logLevelEnvKey    = "MAILER_LOG_LEVEL"
	portEnvKey        = "MAILER_PORT"
	mailerTokenEnvKey = "MAILER_API_KEY"
)

type Config struct {
	MailerToken string
	LogLevel    string
	Port        string
}

func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

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
