package misc

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetHealth(_ context.Context) error {
	return nil
}
