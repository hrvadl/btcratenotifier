package ratesender

import (
	"context"
	"fmt"

	"github.com/hrvadl/btcratenotifier/gw/internal/storage/subscriber"
)

const operation = "ratesender service"

func NewService(rr RecipientRepo, rg RateGetter, s Sender) *Service {
	return &Service{
		repo:       rr,
		rateGetter: rg,
		sender:     s,
	}
}

type RecipientFinder interface {
	FindAll(ctx context.Context) ([]subscriber.Subscriber, error)
}

type RecipientSaver interface {
	Save(ctx context.Context, s subscriber.Subscriber) (int64, error)
}

type RecipientRepo interface {
	RecipientFinder
	RecipientSaver
}

type Sender interface {
	Send(ctx context.Context, html string, emails ...string) error
}

type RateGetter interface {
	GetRate(ctx context.Context) (float32, error)
}

type Service struct {
	repo       RecipientRepo
	sender     Sender
	rateGetter RateGetter
}

func (s *Service) Subscribe(ctx context.Context, mail string) error {
	if _, err := s.repo.Save(ctx, subscriber.Subscriber{Email: mail}); err != nil {
		return fmt.Errorf("%s: failed to save recipient: %w", operation, err)
	}

	return nil
}

func (s *Service) SendToAll(ctx context.Context) error {
	subscribers, err := s.repo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get mails: %w", operation, err)
	}

	r, err := s.rateGetter.GetRate(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get rate: %w", operation, err)
	}

	return s.sender.Send(ctx, fmt.Sprint(r), getMails(subscribers)...)
}

func getMails(s []subscriber.Subscriber) []string {
	mails := make([]string, 0, len(s))
	for _, ss := range s {
		mails = append(mails, ss.Email)
	}
	return mails
}
