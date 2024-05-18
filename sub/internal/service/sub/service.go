package sub

import (
	"context"
	"errors"
	"fmt"

	"github.com/hrvadl/converter/sub/internal/storage/subscriber"
)

const operation = "ratesender service"

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

type Service struct {
	repo      RecipientSaver
	validator Validator
}

func (s *Service) Subscribe(ctx context.Context, mail string) (int64, error) {
	if !s.validator.Validate(mail) {
		return 0, errors.New("invalid email")
	}

	resp, err := s.repo.Save(ctx, subscriber.Subscriber{Email: mail})
	if err != nil {
		return 0, fmt.Errorf("%s: failed to save recipient: %w", operation, err)
	}

	return resp, nil
}
