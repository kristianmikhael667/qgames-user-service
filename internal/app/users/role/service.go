package role

import (
	"context"
	"errors"
	dto "main/internal/dto/users_req_res"
	"main/internal/factory"
	repository "main/internal/repository/user_repo"
	"main/package/util/response"
)

type service struct {
	RoleRepository repository.Role
}

type Service interface {
	CreateRole(ctx context.Context, payload *dto.RoleRequestBody) (*dto.RoleResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		RoleRepository: f.RoleRepository,
	}
}

func (s *service) CreateRole(ctx context.Context, payload *dto.RoleRequestBody) (*dto.RoleResponse, error) {
	var result *dto.RoleResponse

	// Check duplicate
	isExistName, err := s.RoleRepository.ExistByName(ctx, &payload.Name)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistName {
		return result, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("Name Already Exists"))
	}

	data, err := s.RoleRepository.Save(ctx, payload)

	result = &dto.RoleResponse{
		UidRole: data.UidRole.String(),
		Name:    data.Name,
		Desc:    data.Desc,
		Data:    data.Data,
		Status:  data.Status,
	}

	return result, nil
}
