package user

import (
	"context"
	dto "main/internal/dto"
	"main/internal/factory"
	repository "main/internal/repository"
	pkgdto "main/package/dto"
	res "main/package/util/response"
)

type service struct {
	UserRepository repository.User
}

type Service interface {
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.UsersResponse], error)
	UpdateUsers(ctx context.Context, payloads *pkgdto.ByUuidUsersRequest, payload *dto.UpdateUsersReqBody) (*dto.UsersResponse, int16, string, error)
	GetUserDetail(ctx context.Context, roles, iduser string) (*dto.UserResponseAll, int16, string, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
	}
}

func (s *service) Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.UsersResponse], error) {
	users, info, err := s.UserRepository.FindAll(ctx, payload, &payload.Pagination)
	if err != nil {
		return nil, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	var data []dto.UsersResponse

	for _, user := range users {
		data = append(data, dto.UsersResponse{
			Fullname: user.Fullname,
			Email:    user.Email,
		})

	}

	result := new(pkgdto.SearchGetResponse[dto.UsersResponse])
	result.Data = data
	result.PaginationInfo = *info

	return result, nil
}

func (s *service) UpdateUsers(ctx context.Context, payloads *pkgdto.ByUuidUsersRequest, payload *dto.UpdateUsersReqBody) (*dto.UsersResponse, int16, string, error) {
	var result *dto.UsersResponse
	// Update
	data, sc, msg, err := s.UserRepository.UpdateAccount(ctx, payloads.Uid, payload)

	if err != nil {
		return result, sc, msg, err
	}

	result = &dto.UsersResponse{
		Uuid:      data.UidUser.String(),
		Fullname:  data.Fullname,
		Phone:     data.Phone,
		Email:     data.Email,
		Address:   data.Address,
		Profile:   data.Profile,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return result, sc, msg, nil
}

func (s *service) GetUserDetail(ctx context.Context, roles, iduser string) (*dto.UserResponseAll, int16, string, error) {
	var user_data *dto.UserResponseAll

	users, sc, msg, err := s.UserRepository.MyAccount(ctx, iduser)
	if err != nil {
		return nil, sc, msg, err
	}

	user_data = &dto.UserResponseAll{
		Data: dto.UsersResponse{
			Uuid:      users.UidUser.String(),
			Fullname:  users.Fullname,
			Phone:     users.Phone,
			Email:     users.Email,
			Address:   users.Address,
			Profile:   users.Profile,
			CreatedAt: users.CreatedAt,
			UpdatedAt: users.UpdatedAt,
		},
		Roles: roles,
	}

	return user_data, sc, msg, nil
}
