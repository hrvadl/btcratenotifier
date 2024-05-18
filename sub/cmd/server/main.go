package main

import (
	"os"

	"github.com/hrvadl/converter/sub/internal/app"
	"github.com/hrvadl/converter/sub/internal/cfg"
	"github.com/hrvadl/converter/sub/pkg/logger"
)

const (
	source = "sub"
)

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
