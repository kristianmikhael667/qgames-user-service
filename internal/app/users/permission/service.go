package permission

import (
	"context"
	"errors"
	dto "main/internal/dto/users_req_res"
	"main/internal/factory"
	repository "main/internal/repository/user_repo"
	"main/package/util/response"
)

type service struct {
	PermissionRepository repository.Permission
}

type Service interface {
	CreatePermission(ctx context.Context, payload *dto.PermissionRequestBody) (*dto.PermissionResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		PermissionRepository: f.PermissionRepository,
	}
}

func (s *service) CreatePermission(ctx context.Context, payload *dto.PermissionRequestBody) (*dto.PermissionResponse, error) {
	var result *dto.PermissionResponse

	// Check duplicate
	isExistName, err := s.PermissionRepository.ExistByNamePermission(ctx, &payload.Name)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistName {
		return result, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("Name Already Exists"))
	}

	data, err := s.PermissionRepository.Save(ctx, payload)

	result = &dto.PermissionResponse{
		UidPermission: data.UidPermission.String(),
		Name:          data.Name,
		Slug:          data.Slug,
	}

	return result, nil
}
