package sub

import (
	"context"
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

// NewService constructs new Service with provided arguments.
// NOTE: neither of arguments can't be nil, or service will panic in
// the future.
func NewService(rr RecipientSaver, vv Validator) *Service {
	return &Service{
		repo:      rr,
		validator: vv,
	}
}

//go:generate mockgen -destination=./mocks/mock_saver.go -package=mocks . RecipientSaver
type RecipientSaver interface {
	Save(ctx context.Context, s subscriber.Subscriber) (int64, error)
}

//go:generate mockgen -destination=./mocks/mock_validator.go -package=mocks . Validator
type Validator interface {
	Validate(mail string) bool
}

// Service is a main structure, responsible for doing checks
// and calling underlying saver to save subscriber if everything is correct.
type Service struct {
	repo      RecipientSaver
	validator Validator
}

// Subscribe method accepts context and subscriber's mail.
// First of all, it validates subscriber's email.
// Then it call underlying repo to save subscriber:
// If OK returns ID of saved subscriber, if not - returns an error.
func (s *Service) Subscribe(ctx context.Context, mail string) (int64, error) {
	if !s.validator.Validate(mail) {
		return 0, ErrInvalidEmail
	}

	resp, err := s.repo.Save(ctx, subscriber.Subscriber{Email: mail})
	if err == nil {
		return resp, nil
	}

	if errors.Is(err, subscriber.ErrAlreadyExists) {
		return 0, ErrAlreadyExists
	}

	return 0, ErrFailedToSave
}
