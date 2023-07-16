package auth

import (
	"context"
	"errors"
	"fmt"
	dto "main/internal/dto/users_req_res"
	"main/internal/factory"
	"main/internal/pkg/util"
	repository "main/internal/repository/user_repo"
	pkgutil "main/package/util"
	"main/package/util/response"
)

type service struct {
	UserRepository   repository.User
	AssignRepository repository.Assign
}

type Service interface {
	RegisterUsers(ctx context.Context, payload *dto.RegisterUsersRequestBody) (*dto.UserWithJWTResponse, error)
	CheckPhone(ctx context.Context, payload *dto.RegisterUsersRequestBody) (bool, error)
	RequestOtp(ctx context.Context, phone *dto.CheckPhoneReqBody) (string, bool, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository:   f.UserRepository,
		AssignRepository: f.AssignRepository,
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
	isExistPhone, err := s.UserRepository.ExistByPhone(ctx, payload.Phone)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistPhone {
		return result, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("Phone Already Exists"))
	}

	// Hash Password
	hashedPassword, err := pkgutil.HashPassword(payload.Password)
	if err != nil {
		return result, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	payload.Password = hashedPassword

	// Hash Pin
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
	isExistPhone, err := s.UserRepository.ExistByPhone(ctx, payload.Phone)
	if err != nil {
		return true, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExistPhone {
		return true, response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Phone Not Found"))
	}
	return false, err
}

func (s *service) RequestOtp(ctx context.Context, phone *dto.CheckPhoneReqBody) (string, bool, error) {
	// var result *dto.UsersResponse

	isExistPhone, err := s.UserRepository.ExistByPhone(ctx, phone.Phone)
	if err != nil {
		return "Error", true, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if !isExistPhone {
		return "Error", true, response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Phone Not Found"))
	}

	// Request Otp
	requestOtp, status, _, err := s.UserRepository.RequestOtp(ctx, phone.Phone)

	if err != nil {
		fmt.Println("error beo")
		return err.Error(), status, err
	}

	// Assign Role user-default Permission no-topup-balance
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "86976d39-8829-40e8-ada1-45e390c244da")

	// Assign Role user-default Permission common-user
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "2ddfdaba-5024-409f-b641-b6424eb8cb8f")

	// Assign Role user-default Permission check-wallet
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "8487b5a5-8f95-4623-ad23-15fdd507f82b")

	// Assign Role user-default Permission topup-wallet
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "9835a6a5-26aa-4cfc-8cff-ef371ad4ba8b")

	// Assign Role user-default Permission topup-wallet
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "9835a6a5-26aa-4cfc-8cff-ef371ad4ba8b")

	// Assign Role user-default Permission list-product
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "cab88e92-77e6-4744-a1c4-b771c9cda9ef")

	// Assign Role user-default Permission create-trx
	err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "20303ce3-a6fe-4463-a5a3-7d5e333c6b69", "0287be90-35f5-4d94-b588-34020041d23a")

	return "An instruction to verify your phone number has been sent to your phone.", status, nil
}
