package main

import (
	"log/slog"
	"os"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/tests/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/tests/internal/cfg"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("source", "tests"),
		slog.Int("pid", os.Getpid()),
	)

	cfg := cfg.Must(cfg.NewFromEnv())
	l.Info("Parsed config")

	app := app.New(*cfg, l)
	app.MustRun()

	l.Info("Finished load testing")
}
