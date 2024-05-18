package cfg

import (
	"fmt"
	"os"
)

const operation = "config parsing"

const (
	subServiceAddrEnvKey = "SUB_ADDR"
	rateWatchAddrEnvKey  = "RATE_WATCH_ADDR"
	logLevelEnvKey       = "GATEWAY_LOG_LEVEL"
	addrEnvKey           = "GATEWAY_ADDR"
)

type Config struct {
	SubAddr         string
	RateWatcherAddr string
	Addr            string
	LogLevel        string
}

func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewFromEnv() (*Config, error) {
	rwAddr := os.Getenv(rateWatchAddrEnvKey)
	if rwAddr == "" {
		return nil, fmt.Errorf("%s: rate watcher addr can't be empty", operation)
	}

	subAddr := os.Getenv(subServiceAddrEnvKey)
	if subAddr == "" {
		return nil, fmt.Errorf("%s: rate watcher addr can't be empty", operation)
	}

	logLevel := os.Getenv(logLevelEnvKey)
	if logLevel == "" {
		return nil, fmt.Errorf("%s: log level can't be empty", operation)
	}

	port := os.Getenv(addrEnvKey)
	if port == "" {
		return nil, fmt.Errorf("%s: port can't be empty", operation)
	}

	return &Config{
		LogLevel:        logLevel,
		Addr:            port,
		RateWatcherAddr: rwAddr,
		SubAddr:         subAddr,
	}, nil
}
