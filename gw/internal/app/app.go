package app

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/hrvadl/btcratenotifier/gw/internal/cfg"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/ratewatcher"
	ssvc "github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/sub"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/http/handlers/rate"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/http/handlers/sub"
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

	r := http.NewServeMux()
	r.HandleFunc("GET /rate", rh.GetRate)
	r.HandleFunc("POST /subscribe", sh.Subscribe)

	return http.ListenAndServe(net.JoinHostPort("", a.cfg.Port), r)
}
