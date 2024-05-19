package main

import (
	"os"

	"github.com/hrvadl/converter/mailer/internal/app"
	"github.com/hrvadl/converter/mailer/internal/cfg"
	"github.com/hrvadl/converter/mailer/pkg/logger"
)

const source = "mailer"

func main() {
	cfg := cfg.Must(cfg.NewFromEnv())
	l := logger.New(os.Stdout, cfg.LogLevel).With(
		"source", source,
		"pid", os.Getpid(),
	)

	l.Info("Successfuly parsed config and initialized logger")
	app := app.New(*cfg, l)
	go app.MustRun()
	app.GracefulStop()
}
