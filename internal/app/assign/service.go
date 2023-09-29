package assign

import (
	"context"
	"main/internal/dto"
	"main/internal/factory"
	"main/internal/repository"
)

type service struct {
	AssignRepository repository.Assign
}

type Service interface {
	EditAssign(ctx context.Context, payload *dto.ReqAssign) (bool, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		AssignRepository: f.AssignRepository,
	}
}

func (s *service) EditAssign(ctx context.Context, payload *dto.ReqAssign) (bool, error) {
	isTrue, err := s.AssignRepository.EditRolesTopup(ctx, payload)

	if err != nil {
		return isTrue, err
	}

	return isTrue, nil
}
