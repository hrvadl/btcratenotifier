package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/hrvadl/btcratenotifier/gw/internal/cfg"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/ratewatcher"
	ssvc "github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/sub"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/http/handlers/rate"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/http/handlers/sub"
	"github.com/hrvadl/btcratenotifier/gw/pkg/logger"
)

const operation = "app init"

func New(cfg cfg.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

type App struct {
	cfg cfg.Config
	log *slog.Logger
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	rw, err := ratewatcher.NewClient(
		a.cfg.RateWatcherAddr,
		a.log.With("source", "rateWatcherClient"),
	)
	if err != nil {
		return fmt.Errorf("%s: failed to initialize ratewatcher client: %w", operation, err)
	}

	subsvc, err := ssvc.NewClient(a.cfg.SubAddr, a.log.With("source", "subClient"))
	if err != nil {
		return fmt.Errorf("%s: failed to init sub service: %w", operation, err)
	}

	sh := sub.NewHandler(subsvc, a.log.With("source", "subHandler"))
	rh := rate.NewHandler(rw, a.log.With("source", "rateHandler"))

	r := chi.NewRouter()
	r.Use(
		middleware.Heartbeat("/health"),
		middleware.Recoverer,
		middleware.Logger,
		middleware.CleanPath,
	)

	r.Get("/rate", rh.GetRate)
	r.With(
		middleware.AllowContentType("application/x-www-form-urlencoded"),
	).Post("/subscribe", sh.Subscribe)

	srv := newServer(
		r,
		net.JoinHostPort("", a.cfg.Port),
		slog.NewLogLogger(a.log.Handler(), logger.MapLevels(a.cfg.LogLevel)),
	)

	return srv.ListenAndServe()
}
