package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/hrvadl/converter/mailer/internal/cfg"
	"github.com/hrvadl/converter/mailer/internal/platform/mail/resend"
	"github.com/hrvadl/converter/mailer/internal/transport/grpc/server/mailer"
	"github.com/hrvadl/converter/mailer/pkg/logger"
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

	mailer.Register(srv, resend.NewClient(a.cfg.MailerToken), a.log.With("source", "mailerSrv"))
	listener, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("%s: failed to listen on port %s: %w", operation, a.cfg.Port, err)
	}

	return srv.Serve(listener)
}
