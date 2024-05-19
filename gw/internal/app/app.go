package app

//nolint:revive
import (
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/hrvadl/converter/gw/docs"
	"github.com/hrvadl/converter/gw/internal/cfg"
	"github.com/hrvadl/converter/gw/internal/transport/grpc/clients/ratewatcher"
	ssvc "github.com/hrvadl/converter/gw/internal/transport/grpc/clients/sub"
	"github.com/hrvadl/converter/gw/internal/transport/http/handlers/rate"
	"github.com/hrvadl/converter/gw/internal/transport/http/handlers/sub"
	"github.com/hrvadl/converter/gw/pkg/logger"
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

	a.log.Info("Starting web server", "addr", a.cfg.Addr)
	srv := newServer(
		r,
		a.cfg.Addr,
		slog.NewLogLogger(a.log.Handler(), logger.MapLevels(a.cfg.LogLevel)),
	)

	return srv.ListenAndServe()
}
