package sender

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hrvadl/converter/sub/internal/storage/subscriber"
)

const (
	operation = "sender cron job"
	subject   = "BTC to UAH rate exchange"
)

func New(
	m Mailer,
	sg SubscriberGetter,
	mf RateMessageFormatter,
	rg RateGetter,
	log *slog.Logger,
) *Service {
	return &Service{
		mailer:     m,
		subGetter:  sg,
		formatter:  mf,
		rateGetter: rg,
		log:        log,
	}
}

//go:generate mockgen -destination=./mocks/mock_rategetter.go -package=mocks . RateGetter
type RateGetter interface {
	GetRate(ctx context.Context) (float32, error)
}

//go:generate mockgen -destination=./mocks/mock_subgetter.go -package=mocks . SubscriberGetter
type SubscriberGetter interface {
	Get(ctx context.Context) ([]subscriber.Subscriber, error)
}

//go:generate mockgen -destination=./mocks/mock_formatter.go -package=mocks . RateMessageFormatter
type RateMessageFormatter interface {
	Format(r float32) string
}

//go:generate mockgen -destination=./mocks/mock_mailer.go -package=mocks . Mailer
type Mailer interface {
	Send(ctx context.Context, msg, subject string, to ...string) error
}

type Service struct {
	mailer     Mailer
	formatter  RateMessageFormatter
	subGetter  SubscriberGetter
	rateGetter RateGetter
	log        *slog.Logger
}

func (w *Service) Send(ctx context.Context) error {
	subs, err := w.subGetter.Get(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get subscribers: %w", operation, err)
	}

	if len(subs) == 0 {
		return fmt.Errorf("%s: can't send emails when subscribers are empty", operation)
	}

	r, err := w.rateGetter.GetRate(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get rate: %w", operation, err)
	}

	return w.mailer.Send(ctx, w.formatter.Format(r), subject, mapSubsToMails(subs)...)
}

func mapSubsToMails(s []subscriber.Subscriber) []string {
	mails := make([]string, 0, len(s))
	for i := range s {
		mails = append(mails, s[i].Email)
	}
	return mails
}
