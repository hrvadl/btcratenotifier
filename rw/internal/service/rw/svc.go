package rw

import "context"

func NewService(source RateSource) *Service {
	return &Service{
		sources: []RateSource{source},
	}
}

//go:generate mockgen -destination=./mocks/mock_source.go -package=mocks . RateSource
type RateSource interface {
	Convert(ctx context.Context) (float32, error)
}

type Service struct {
	sources []RateSource
}

func (s *Service) Convert(ctx context.Context) (float32, error) {
	var (
		rate float32
		err  error
	)

	for _, svc := range s.sources {
		rate, err = svc.Convert(ctx)
		if err == nil {
			return rate, nil
		}
	}

	return rate, err
}

func (s *Service) SetNext(source ...RateSource) {
	s.sources = append(s.sources, source...)
}
