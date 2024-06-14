package sender

import (
	"context"
	"log/slog"
	"time"
)

// NewCronJobAdapter constructs CronJobAdapter for Sender interface
// compatible structure.
// NOTE: neither of arguments can't be nil, or service will panic in the future.
func NewCronJobAdapter(s Sender, timeout time.Duration, log *slog.Logger) *CronJobAdapter {
	return &CronJobAdapter{
		sender:  s,
		log:     log,
		timeout: timeout,
	}
}

//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks . Sender
type Sender interface {
	Send(ctx context.Context) error
}

// CronJobAdapter is a handy wrapper to help Sender compatible
// structure fit to the CronJob required interface.
type CronJobAdapter struct {
	timeout time.Duration
	sender  Sender
	log     *slog.Logger
}

// Do method log's each call then creates context with default timeout of 10 seconds
// and then executes original function, returning the error if any.
func (c *CronJobAdapter) Do() error {
	c.log.Info("Sending mails in cron job")
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.sender.Send(ctx)
}
