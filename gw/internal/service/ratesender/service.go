package ratesender

import (
	"context"
	"fmt"

	"github.com/hrvadl/btcratenotifier/gw/internal/transport/grpc/clients/mailer"
)

const operation = "ratesender service"

func NewService() *Service {
	return &Service{}
}

type RecipientFinder interface {
	Find(ctx context.Context, mail string) (string, error)
	FindAll(ctx context.Context) ([]string, error)
}

type RecipientSaver interface {
	Save(ctx context.Context, mail string) error
}

type RecipientRepo interface {
	RecipientFinder
	RecipientSaver
}

type Sender interface {
	Send(ctx context.Context, emails ...mailer.SendOptions) error
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
	if find, _ := s.repo.Find(ctx, mail); find != "" {
		return fmt.Errorf("%s: cannot add recipient when it already exists", operation)
	}

	if err := s.repo.Save(ctx, mail); err != nil {
		return fmt.Errorf("%s: failed to save recipient: %w", operation, err)
	}

	return nil
}

func (s *Service) SendToAll(ctx context.Context) error {
	mails, err := s.repo.FindAll(ctx)
	if err != nil {
		fmt.Errorf("%s: failed to get mails: %w", operation, err)
	}

	r, err := s.rateGetter.GetRate(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get rate: %w", operation, err)
	}

	return s.sender.Send(ctx, mailsToMailerOptions(r, mails)...)
}

func mailsToMailerOptions(rate float32, mails []string) []mailer.SendOptions {
	opt := make([]mailer.SendOptions, 0, len(mails))
	for _, m := range mails {
		opt = append(opt, mailer.SendOptions{
			To:      m,
			Payload: fmt.Sprint(rate),
		})
	}
	return opt
}
