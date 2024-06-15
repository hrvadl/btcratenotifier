package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/tests/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/tests/internal/prom"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/tests/internal/tests/gw"
)

// New constructs new App with provided arguments.
// NOTE: than neither cfg or log can't be nil or App will panic.
func New(cfg cfg.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

// App is a thin abstraction used to initialize all the dependencies,
// db connections, and GRPC server/clients. Could return an error if any
// of described above steps failed.
type App struct {
	cfg cfg.Config
	log *slog.Logger
}

// MustRun is a wrapper around App.Run() function which could be handly
// when it's called from the main goroutine and we don't need to handler
// an error.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run method creates new prometheus server then initializes all load tests,
// after that runs all tests and finally
// starts exposing metrics on the provided port. Could return an error if any of
// described above steps failed
func (a *App) Run() error {
	srv, pm, err := prom.NewServer(a.cfg)
	if err != nil {
		return err
	}

	go func() {
		a.log.Info("Starting the server...", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil {
			a.log.Info("Error while running server...", slog.Any("err", err))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	defer func() {
		a.log.Info("Starting down the server...", slog.String("addr", srv.Addr))
		if err := srv.Shutdown(ctx); err != nil {
			a.log.Info("Failed to shut down server...", slog.Any("err", err))
		}
	}()

	start := time.Now()
	a.log.Info("Starting tests")
	lt := gw.NewLoadTest(pm, a.cfg.GatewayAddr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// lt.GetRate()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lt.Subscribe()
	}()

	wg.Wait()
	a.log.Info("Finished tests", slog.Any("took", time.Since(start).Seconds()))
	return nil
}
