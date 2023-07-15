package auth

import (
	"context"
	"errors"
	dto "main/internal/dto/users_req_res"
	"main/internal/factory"
	"main/internal/pkg/util"
	repository "main/internal/repository/user_repo"
	pkgo "main/package/dto"
	pkgutil "main/package/util"
	"main/package/util/response"
)

type service struct {
	UserRepository repository.User
}

type Service interface {
	RegisterUsers(ctx context.Context, payload *dto.RegisterUsersRequestBody) (*dto.UserWithJWTResponse, error)
	CheckPhone(ctx context.Context, payload *dto.RegisterUsersRequestBody) (bool, error)
	CheckPhonesLogin(ctx context.Context, phone *pkgo.ByPhoneNumber) (bool, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
	}
}

func (s *service) RegisterUsers(ctx context.Context, payload *dto.RegisterUsersRequestBody) (*dto.UserWithJWTResponse, error) {
	var result *dto.UserWithJWTResponse

	// Check Email
	isExistEmail, err := s.UserRepository.ExistByEmail(ctx, &payload.Email)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistEmail {
		return result, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("Email Already Exists"))
	}

	// Check Phone
	isExistPhone, err := s.UserRepository.ExistByPhone(ctx, &payload.Phone)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistPhone {
		return result, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("Phone Already Exists"))
	}

	hashedPassword, err := pkgutil.HashPassword(payload.Password)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	payload.Password = hashedPassword

	data, err := s.UserRepository.Save(ctx, payload)

	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	claims := util.CreateJWTClaims(data.UidUser.String(), data.Email, data.Phone)
	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("Error generate token"),
		)
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      data.UidUser.String(),
			Fullname:  data.Fullname,
			Phone:     data.Phone,
			Email:     data.Email,
			Address:   data.Address,
			Profile:   data.Profile,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
		},
		Token: token,
	}

	return result, nil
}

func (s *service) CheckPhone(ctx context.Context, payload *dto.RegisterUsersRequestBody) (bool, error) {
	// Check Phone
	isExistPhone, err := s.UserRepository.ExistByPhone(ctx, &payload.Phone)
	if err != nil {
		return true, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistPhone {
		return true, response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Phone Not Found"))
	}
	return false, err
}

func (s *service) CheckPhonesLogin(ctx context.Context, phone *pkgo.ByPhoneNumber) (bool, error) {
	isExistPhone, err := s.UserRepository.ExistByPhone(ctx, &phone.Phone)
	if err != nil {
		return true, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if !isExistPhone {
		return true, response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Phone Not Found"))
	}
	return false, err
}
