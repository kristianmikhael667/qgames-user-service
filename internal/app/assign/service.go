package assign

import (
	"context"
	"main/internal/dto"
	"main/internal/factory"
	"main/internal/repository"

	"github.com/labstack/echo/v4"
)

type service struct {
	AssignRepository repository.Assign
}

type Service interface {
	EditAssign(c echo.Context, ctx context.Context, payload *dto.ReqAssign) (bool, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		AssignRepository: f.AssignRepository,
	}
}

func (s *service) EditAssign(c echo.Context, ctx context.Context, payload *dto.ReqAssign) (bool, error) {
	isTrue, err := s.AssignRepository.EditRolesTopup(c, ctx, payload)

	if err != nil {
		return isTrue, err
	}

	return isTrue, nil
}
