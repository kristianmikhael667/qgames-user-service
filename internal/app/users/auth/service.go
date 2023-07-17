package auth

import (
	"context"
	"errors"
	"main/helper"
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
	VerifyOtp(ctx context.Context, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, string, error)
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
	// Request Otp
	requestOtp, status, _, err := s.UserRepository.RequestOtp(ctx, phone.Phone)

	if err != nil {
		return err.Error(), status, err
	}

	if status == true {
		// Assign Role user-default Permission no-topup-balance
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "no-topup-balance")

		// Assign Role user-default Permission common-user
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "common-user")

		// Assign Role user-default Permission check-wallet
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "check-wallet")

		// Assign Role user-default Permission topup-wallet
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "topup-wallet")

		// Assign Role user-default Permission list-product
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "list-product")

		// Assign Role user-default Permission create-trx
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "create-trx")
	}

	return "An instruction to verify your phone number has been sent to your phone.", status, nil
}

func (s *service) VerifyOtp(ctx context.Context, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, string, error) {
	var result *dto.UserWithJWTResponse

	// Check Email
	response, verifyOtp, msg, err := s.UserRepository.VerifyOtp(ctx, validotp.Phone, validotp.Otp)
	if err != nil {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, err
	}
	if verifyOtp == false {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, err
	}

	claims := util.CreateJWTClaims(response.UidUser.String(), response.Email, response.Phone)
	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, msg, err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      response.UidUser.String(),
			Fullname:  response.Fullname,
			Phone:     response.Phone,
			Email:     response.Email,
			Address:   response.Address,
			Profile:   response.Profile,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		},
		Token: token,
	}
	return result, msg, nil
}
