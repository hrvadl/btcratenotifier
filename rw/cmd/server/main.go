package main

import (
	"log/slog"
	"os"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/kit/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/cfg"
)

const source = "rateWatcher"

func main() {
	cfg := cfg.Must(cfg.NewFromEnv())
	l := logger.New(os.Stdout, cfg.LogLevel).With(
		slog.String("source", source),
		slog.Int("pid", os.Getpid()),
	)

	l.Info("Successfully parsed config and initialized logger")
	app := app.New(*cfg, l)
	go app.MustRun()
	app.GracefulStop()
}
