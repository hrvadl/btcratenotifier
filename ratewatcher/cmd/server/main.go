package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hrvadl/btcratenotifier/pkg/logger"

	"github.com/hrvadl/ratewatcher/internal/cfg"
	"github.com/hrvadl/ratewatcher/internal/platform/rates/cryptocompare"
)

const source = "rateWatcher"

func main() {
	cfg := cfg.Must(cfg.NewFromEnv())
	l := logger.New(os.Stdout, cfg.LogLevel).With(
		"source", source,
		"pid", os.Getpid(),
	)

	l.Info("Successfuly parsed config...")
	rw := cryptocompare.NewClient(cfg.ExchangeServiceToken, cfg.ExchangeServiceBaseURL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rate, err := rw.BTCToUAH(ctx)
	if err != nil {
		panic(err)
	}

	l.Info("Got a fresh rate!", "rate", fmt.Sprintf("%.2f", rate))
}
