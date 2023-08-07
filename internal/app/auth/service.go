package auth

import (
	"context"
	"errors"
	"fmt"
	"main/helper"
	dto "main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
	repository "main/internal/repository"
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
	RequestOtp(ctx context.Context, phone *dto.CheckPhoneReqBody) (string, int, bool, error)
	VerifyOtp(ctx context.Context, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, string, int16, error)
	LoginPin(ctx context.Context, loginotp *dto.LoginByPin) (*dto.UserWithJWTResponse, string, int16, error)
	LoginAdmin(ctx context.Context, loginadmin *dto.LoginAdmin) (*dto.UserWithJWTResponse, string, int, error)
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

	claims := util.CreateJWTClaims(data.UidUser.String(), data.Email, data.Phone, "nil", nil, false)
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

func (s *service) RequestOtp(ctx context.Context, phone *dto.CheckPhoneReqBody) (string, int, bool, error) {
	// Request Otp
	requestOtp, sc, status, msg, err := s.UserRepository.RequestOtp(ctx, phone)

	if err != nil {
		return err.Error(), sc, status, err
	}

	if status == true {
		// Assign Role user-default and Permission
		err = s.AssignRepository.Assign(ctx, requestOtp.UidUser.String(), "user-default", "common-user,no-topup-balance,check-wallet,topup-wallet,list-product,create-trx")
		if err != nil {
			return "error assign", sc, status, err
		}
	}

	return msg, sc, status, nil
}

func (s *service) VerifyOtp(ctx context.Context, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, string, int16, error) {
	var result *dto.UserWithJWTResponse
	// Check Email
	response, verifyOtp, msg, err := s.UserRepository.VerifyOtp(ctx, validotp.Phone, validotp.Otp)
	if err != nil {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, 403, err
	}
	if verifyOtp == false {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, 401, err
	}

	// Get all assign and loop
	response_assign, err := s.UserRepository.GetAssignUsers(ctx, response.UidUser.String())

	if err != nil {
		helper.Logger("error", "Error get assign user service", "Rc: "+string(rune(403)))
	}

	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(response.UidUser.String(), response.Email, response.Phone, firstRole, permissions, false)

	// Update Limit
	statuscode, msg, err := s.UserRepository.UpdateAttemptOtp(ctx, validotp.Phone)
	if statuscode != 201 {
		return result, msg, statuscode, err
	}

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, msg, 401, err
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
	return result, msg, 201, nil
}

func (s *service) LoginPin(ctx context.Context, loginotp *dto.LoginByPin) (*dto.UserWithJWTResponse, string, int16, error) {
	var result *dto.UserWithJWTResponse

	// Login Pin
	responses, sc, msg, err := s.UserRepository.LoginByPin(ctx, loginotp)

	if err != nil {
		return result, msg, sc, err
	}
	fmt.Println("ssss ", responses.Phone)
	// Get all assign and loop
	response_assign, err := s.UserRepository.GetAssignUsers(ctx, responses.UidUser.String())

	if err != nil {
		helper.Logger("error", "Error get assign user service", "Rc: "+string(rune(403)))
	}

	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(responses.UidUser.String(), responses.Email, responses.Phone, firstRole, permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, msg, 401, err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      responses.UidUser.String(),
			Fullname:  responses.Fullname,
			Phone:     responses.Phone,
			Email:     responses.Email,
			Address:   responses.Address,
			Profile:   responses.Profile,
			CreatedAt: responses.CreatedAt,
			UpdatedAt: responses.UpdatedAt,
		},
		Token: token,
	}
	return result, msg, 201, nil
}

func (s *service) LoginAdmin(ctx context.Context, loginadmin *dto.LoginAdmin) (*dto.UserWithJWTResponse, string, int, error) {
	var result *dto.UserWithJWTResponse

	// Login Admin
	responses, sc, msg, err := s.UserRepository.LoginAdmin(ctx, loginadmin)

	if err != nil {
		return result, msg, sc, err
	}

	claims := util.CreateJWTClaims(responses.UidUser.String(), responses.Email, responses.Phone, "nil", nil, true)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, msg, 401, err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      responses.UidUser.String(),
			Fullname:  responses.Fullname,
			Phone:     responses.Phone,
			Email:     responses.Email,
			Address:   responses.Address,
			Profile:   responses.Profile,
			CreatedAt: responses.CreatedAt,
			UpdatedAt: responses.UpdatedAt,
		},
		Token: token,
		Admin: true,
	}
	return result, msg, sc, nil
}
