package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/cfg"
	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/platform/rates/exchangerate"
	"github.com/hrvadl/btcratenotifier/ratewatcher/internal/transport/grpc/server/ratewatcher"
	"github.com/hrvadl/btcratenotifier/ratewatcher/pkg/logger"
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
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logger.NewServerGRPCMiddleware(a.log),
	))

	ratewatcher.Register(
		srv,
		exchangerate.NewClient(a.cfg.ExchangeServiceToken, a.cfg.ExchangeServiceBaseURL),
		a.log.With("source", "rateWatcherSrv"),
	)
	a.log.Info("Successfuly initialized all deps")

	listener, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("%s: failed to listen on tcp port %s: %w", operation, a.cfg.Port, err)
	}

	return srv.Serve(listener)
}
