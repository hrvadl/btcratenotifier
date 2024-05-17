package cfg

import (
	"fmt"
	"os"
)

const operation = "config parsing"

const (
	mailerServiceAddrEnvKey = "MAILER_ADDR"
	rateWatchAddrEnvKey     = "RATE_WATCH_ADDR"
	logLevelEnvKey          = "GATEWAY_LOG_LEVEL"
	portEnvKey              = "GATEWAY_PORT"
	dsnEnvKey               = "GATEWAY_DSN"
)

type Config struct {
	MailerAddr      string
	Dsn             string
	RateWatcherAddr string
	Port            string
	LogLevel        string
}

func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewFromEnv() (*Config, error) {
	mAddr := os.Getenv(mailerServiceAddrEnvKey)
	if mAddr == "" {
		return nil, fmt.Errorf("%s: mailer addr can't be empty", operation)
	}

	rwAddr := os.Getenv(rateWatchAddrEnvKey)
	if rwAddr == "" {
		return nil, fmt.Errorf("%s: rate watcher addr can't be empty", operation)
	}

	logLevel := os.Getenv(logLevelEnvKey)
	if logLevel == "" {
		return nil, fmt.Errorf("%s: log level can't be empty", operation)
	}

	port := os.Getenv(portEnvKey)
	if port == "" {
		return nil, fmt.Errorf("%s: port can't be empty", operation)
	}

	dsn := os.Getenv(dsnEnvKey)
	if dsn == "" {
		return nil, fmt.Errorf("%s: dsn can't be empty", operation)
	}

	return &Config{
		LogLevel:        logLevel,
		Port:            port,
		RateWatcherAddr: rwAddr,
		MailerAddr:      mAddr,
		Dsn:             dsn,
	}, nil
}
