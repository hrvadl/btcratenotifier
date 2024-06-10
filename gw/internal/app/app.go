package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/docs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/grpc/clients/ratewatcher"
	ssvc "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/grpc/clients/sub"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/http/handlers/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/http/handlers/sub"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/pkg/logger"
)

const operation = "app init"

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

// Run method creates new GRPC server then initializes MySQL DB connection,
// after that initializes all necessary domain related services and finally
// starts listening on the provided ports. Could return an error if any of
// described above steps failed
func (a *App) Run() error {
	rw, err := ratewatcher.NewClient(
		a.cfg.RateWatcherAddr,
		a.log.With(slog.String("source", "rateWatcherClient")),
	)
	if err != nil {
		return fmt.Errorf("%s: failed to initialize ratewatcher client: %w", operation, err)
	}

	subsvc, err := ssvc.NewClient(a.cfg.SubAddr, a.log.With(slog.String("source", "subClient")))
	if err != nil {
		return fmt.Errorf("%s: failed to init sub service: %w", operation, err)
	}

	sh := sub.NewHandler(subsvc, a.log.With(slog.String("source", "subHandler")))
	rh := rate.NewHandler(rw, a.log.With(slog.String("source", "rateHandler")))

	r := chi.NewRouter()
	r.Use(
		middleware.Heartbeat("/health"),
		middleware.Recoverer,
		middleware.Logger,
		middleware.CleanPath,
		middleware.SetHeader("Content-Type", "application/octet-stream"),
	)

	r.Route("/api", func(r chi.Router) {
		r.Get("/rate", rh.GetRate)
		r.With(
			middleware.AllowContentType("application/x-www-form-urlencoded"),
		).Post("/subscribe", sh.Subscribe)
	})

	if a.cfg.LogLevel == "DEBUG" {
		r.Get("/docs/*", httpSwagger.WrapHandler)
	}

	a.log.Info("Starting web server", slog.String("addr", a.cfg.Addr))
	srv := newServer(
		r,
		a.cfg.Addr,
		slog.NewLogLogger(a.log.Handler(), logger.MapLevels(a.cfg.LogLevel)),
	)

	return srv.ListenAndServe()
}

// GracefulStop method gracefully stop the server. It listens to the OS sigals.
// After it receives signal it terminates all currently active servers,
// client, connections (if any) and gracefully exits.
func (a *App) GracefulStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	signal := <-ch
	a.log.Info("Received stop signal. Terminating...", slog.Any("signal", signal))
	a.log.Info("Successfully terminated server. Bye!")
}
