package main

import (
	"os"

	"github.com/hrvadl/btcratenotifier/gw/internal/app"
	"github.com/hrvadl/btcratenotifier/gw/internal/cfg"
	"github.com/hrvadl/btcratenotifier/gw/pkg/logger"
)

const source = "gateway"

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
