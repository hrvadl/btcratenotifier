package sender

import (
	"context"
	"log/slog"
	"time"
)

func NewCronJobAdapter(s Sender, log *slog.Logger) *CronJobAdapter {
	return &CronJobAdapter{
		sender: s,
		log:    log,
	}
}

//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks . Sender
type Sender interface {
	Send(ctx context.Context) error
}

type CronJobAdapter struct {
	sender Sender
	log    *slog.Logger
}

func (c *CronJobAdapter) Do() error {
	c.log.Info("Sending mails in cron job")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return c.sender.Send(ctx)
}
