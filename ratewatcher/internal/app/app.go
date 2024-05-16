package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/cfg"
	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/platform/rates/cryptocompare"
	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/transport/grpc/server/ratewatcher"
)

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
	srv := grpc.NewServer()

	ratewatcher.Register(
		srv,
		cryptocompare.NewClient(a.cfg.ExchangeServiceToken, a.cfg.ExchangeServiceBaseURL),
		a.log.With("source", "rateWatcherSrv"),
	)
	a.log.Info("Successfuly initialized all deps")

	listener, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on tcp port %s: %w", a.cfg.Port, err)
	}

	return srv.Serve(listener)
}
