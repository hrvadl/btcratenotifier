package cooldown

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/hrvadl/btcratenotifier/sub/internal/storage/date"
)

const operation = "cooled down sender"

func NewSenderDecorator(
	interval time.Duration,
	lsr LastSentRepo,
	s Sender,
	log *slog.Logger,
) *Service {
	return &Service{
		interval: interval,
		repo:     lsr,
		sender:   s,
		log:      log,
	}
}

type Sender interface {
	Send(context.Context) error
}

type LastSentRepo interface {
	UpdateLatestSent(ctx context.Context, d time.Time) error
	GetLatestSent(ctx context.Context) (*date.LatestSent, error)
}

type Service struct {
	interval time.Duration
	repo     LastSentRepo
	sender   Sender
	log      *slog.Logger
}

func (d *Service) Send(ctx context.Context) error {
	skip, err := d.needCooldown(ctx)
	if err != nil {
		return err
	}

	if skip {
		d.log.Info("Cooling down... Skiping sending cuz it was already sent!")
		return nil
	}

	if err := d.sender.Send(ctx); err != nil {
		return err
	}

	return d.updateSentDate(ctx)
}

func (d *Service) needCooldown(ctx context.Context) (bool, error) {
	date, err := d.repo.GetLatestSent(ctx)
	if err != nil {
		return false, fmt.Errorf("%s: failed to get latest sent date: %w", operation, err)
	}

	return time.Since(date.Date).Seconds() < d.interval.Seconds(), nil
}

func (d *Service) updateSentDate(ctx context.Context) error {
	return d.repo.UpdateLatestSent(ctx, time.Now())
}
