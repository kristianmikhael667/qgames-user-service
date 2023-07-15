package user

import (
	"context"
	dto "main/internal/dto/users_req_res"
	"main/internal/factory"
	repository "main/internal/repository/user_repo"
	pkgdto "main/package/dto"
	res "main/package/util/response"
)

type service struct {
	UserRepository repository.User
}

type Service interface {
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.UsersResponse], error)
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
