package main

import (
	"os"

	"github.com/hrvadl/btcratenotifier/pkg/logger"

	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/app"
	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/cfg"
)

const source = "rateWatcher"

func main() {
	cfg := cfg.Must(cfg.NewFromEnv())
	l := logger.New(os.Stdout, cfg.LogLevel).With(
		"source", source,
		"pid", os.Getpid(),
	)

	l.Info("Successfuly parsed config and initialized logger")
	app := app.New(*cfg, l)
	app.MustRun()
}
