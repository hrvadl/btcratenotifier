package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/hrvadl/converter/rw/internal/cfg"
	"github.com/hrvadl/converter/rw/internal/platform/rates/exchangerate"
	"github.com/hrvadl/converter/rw/internal/transport/grpc/server/ratewatcher"
	"github.com/hrvadl/converter/rw/pkg/logger"
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
// after that initializes all neccessary domain related services and finally
// starts listening on the provided ports. Could return an error if any of
// described above steps failed
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
