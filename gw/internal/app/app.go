package app

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/hrvadl/btcratenotifier/gw/internal/cfg"
	"github.com/hrvadl/btcratenotifier/gw/internal/service/ratesender"
	"github.com/hrvadl/btcratenotifier/gw/internal/storage/platform/db"
	"github.com/hrvadl/btcratenotifier/gw/internal/storage/subscriber"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/mailer"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/ratewatcher"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/http/handlers/rate"
	"github.com/hrvadl/btcratenotifier/gw/internal/transport/http/handlers/sender"
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

	mc, err := mailer.NewClient(a.cfg.MailerAddr, a.log.With("source", "mailerClient"))
	if err != nil {
		return fmt.Errorf("%s: failed to initialize mailer client: %w", operation, err)
	}

	db := db.Must(db.NewConn(a.cfg.Dsn))
	sr := subscriber.NewRepo(db)

	rss := ratesender.NewService(sr, rw, mc)
	sh := sender.NewHandler(rss, a.log)

	rh := rate.NewHandler(rw, a.log)

	r := http.NewServeMux()
	r.HandleFunc("POST /subscribe", sh.Subscribe)
	r.HandleFunc("POST /sendEmail", sh.SendToAll)
	r.HandleFunc("GET /rate", rh.GetRate)

	return http.ListenAndServe(net.JoinHostPort("", a.cfg.Port), r)
}
