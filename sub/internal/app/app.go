package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/hrvadl/btcratenotifier/sub/internal/cfg"
	"github.com/hrvadl/btcratenotifier/sub/internal/service/cron"
	"github.com/hrvadl/btcratenotifier/sub/internal/service/sender"
	"github.com/hrvadl/btcratenotifier/sub/internal/service/sender/formatter"
	subs "github.com/hrvadl/btcratenotifier/sub/internal/service/sub"
	"github.com/hrvadl/btcratenotifier/sub/internal/storage/platform/db"
	"github.com/hrvadl/btcratenotifier/sub/internal/storage/subscriber"
	"github.com/hrvadl/btcratenotifier/sub/internal/transport/grpc/clients/mailer"
	"github.com/hrvadl/btcratenotifier/sub/internal/transport/grpc/clients/ratewatcher"
	"github.com/hrvadl/btcratenotifier/sub/internal/transport/grpc/server/sub"
	"github.com/hrvadl/btcratenotifier/sub/pkg/logger"
)

const operation = "app init"

const (
	cronJobHour   = 12
	cronJobMinute = 0
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
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logger.NewServerGRPCMiddleware(a.log),
	))

	db, err := db.NewConn(a.cfg.Dsn)
	if err != nil {
		return fmt.Errorf("%s: failed to init db: %w", operation, err)
	}

	sr := subscriber.NewRepo(db)
	svc := subs.NewService(sr)
	sub.Register(srv, svc, a.log.With("source", "sub"))

	m, err := mailer.NewClient(a.cfg.MailerAddr, a.cfg.MailerFromAddr, a.log)
	if err != nil {
		return fmt.Errorf("%s: failed to connect to mailer service: %w", operation, err)
	}

	sg := subscriber.NewRepo(db)
	fmter := formatter.NewWithDate()
	rw, err := ratewatcher.NewClient(a.cfg.RateWatcherAddr, a.log.With("source", "rateWatcher"))
	if err != nil {
		return fmt.Errorf("%s: failed to connect to rate watcher: %w", operation, err)
	}

	mailSender := sender.New(
		m,
		sg,
		fmter,
		rw,
		a.log.With("source", "cron sender"),
	)

	cronAdapter := sender.NewCronJobAdapter(mailSender, a.log.With("source", "adapter"))
	job := cron.NewDailyJob(cronJobHour, cronJobMinute, a.log.With("source", "cron"))
	job.Do(cronAdapter.Do)

	l, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("%s: failed to start listener on port %s: %w", operation, a.cfg.Port, err)
	}

	return srv.Serve(l)
}
