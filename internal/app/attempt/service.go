package attempt

import (
	"context"
	dto "main/internal/dto"
	"main/internal/factory"
	repository "main/internal/repository"
)

type service struct {
	AttemptRepository repository.Attempt
}

type Service interface {
	ResetOtpService(ctx context.Context, payload *dto.RequestReset) (string, int, error)
	ResetPinService(ctx context.Context, payload *dto.RequestReset) (string, int, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		AttemptRepository: f.AttemptRepository,
	}
}

func (s *service) ResetOtpService(ctx context.Context, payload *dto.RequestReset) (string, int, error) {
	sc, msg, err := s.AttemptRepository.ResetAttemptOtp(ctx, payload)
	if err != nil {
		return msg, sc, err
	}
	return msg, sc, nil
}

func (s *service) ResetPinService(ctx context.Context, payload *dto.RequestReset) (string, int, error) {
	sc, msg, err := s.AttemptRepository.ResetAttemptPin(ctx, payload)
	if err != nil {
		return msg, sc, err
	}
	return msg, sc, nil
}
