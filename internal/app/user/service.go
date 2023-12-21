package user

import (
	"context"
	"fmt"
	dto "main/internal/dto"
	"main/internal/factory"
	"main/internal/pkg/util"
	repository "main/internal/repository"
	"main/package/constant"
	pkgdto "main/package/dto"
	res "main/package/util/response"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type service struct {
	UserRepository     repository.User
	AssignRepository   repository.Assign
	SessionRepository  repository.Session
	FcmTokenRepository repository.Fcmtoken
}

type Service interface {
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.UsersResponse], error)
	FindIdUser(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.UsersResponse, error)
	UpdateUsers(ctx context.Context, payload *dto.UpdateUsersReqBody) (*dto.UserWithJWTResponse, int, string, error)
	GetUserDetail(c echo.Context, ctx context.Context, roles, iduser string) (*dto.UsersResponse, int, string, error)
	ResetPin(ctx context.Context, roles, uiduser string, payload *dto.ConfirmPin) (*dto.UsersResponse, int, string, error)
	Logout(c echo.Context, ctx context.Context, uiduser string) (string, int, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository:     f.UserRepository,
		AssignRepository:   f.AssignRepository,
		SessionRepository:  f.SessionRepository,
		FcmTokenRepository: f.FcmTokenRepository,
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

func (s *service) FindIdUser(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.UsersResponse, error) {
	var result dto.UsersResponse

	// Find Users
	users, err := s.UserRepository.FindIDUser(ctx, payload.ID)
	if err != nil {

		if err == constant.RECORD_NOT_FOUND {
			return nil, err
		}
		return nil, err
	}

	// Find Roles
	assigns, err := s.AssignRepository.FindUserID(ctx, payload.ID)
	if err != nil {

		if err == constant.RECORD_NOT_FOUND {
			return nil, err
		}
		return nil, err
	}

	result.Uuid = users.UidUser.String()
	result.Fullname = users.Fullname
	result.Phone = users.Phone
	result.Email = users.Email
	result.Address = users.Address
	result.Profile = users.Profile
	result.CreatedAt = users.CreatedAt
	result.UpdatedAt = users.UpdatedAt
	result.Roles = assigns.Roles

	return &result, err
}

func (s *service) UpdateUsers(ctx context.Context, payload *dto.UpdateUsersReqBody) (*dto.UserWithJWTResponse, int, string, error) {
	var result *dto.UserWithJWTResponse
	// Update
	data, sc, msg, err := s.UserRepository.UpdateAccount(ctx, payload)

	if err != nil {
		return result, sc, msg, err
	}

	// Get all assign and loop
	response_assign, err := s.AssignRepository.GetAssignUsers(ctx, data.UidUser.String())

	if err != nil {
		log.Print("Error get assign user service ", 403)
	}
	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(data.UidUser.String(), data.Email, data.Phone, firstRole, permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, 401, msg, err
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

	return result, sc, msg, nil
}

func (s *service) GetUserDetail(c echo.Context, ctx context.Context, roles, iduser string) (*dto.UsersResponse, int, string, error) {
	var user_data *dto.UsersResponse

	users, sc, msg, err := s.UserRepository.MyAccount(ctx, iduser)
	if err != nil {
		return nil, sc, msg, err
	}

	msgSess, scSess, _, errSess := s.SessionRepository.CheckSession(c, ctx, users.UidUser.String(), users.Phone, sc, msg)
	if errSess != nil {
		return nil, scSess, msgSess, err
	}
	fmt.Println("msk ini sc ", scSess)
	if scSess == 403 || scSess == 201 {
		return nil, scSess, msgSess, err
	}

	user_data = &dto.UsersResponse{
		Uuid:      users.UidUser.String(),
		Fullname:  users.Fullname,
		Phone:     users.Phone,
		Email:     users.Email,
		Address:   users.Address,
		Profile:   users.Profile,
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
		Roles:     roles,
	}

	return user_data, sc, msg, nil
}

func (s *service) ResetPin(ctx context.Context, uiduser string, roles string, payload *dto.ConfirmPin) (*dto.UsersResponse, int, string, error) {
	var result *dto.UsersResponse
	// Reset PIN
	data, sc, msg, err := s.UserRepository.ResetPin(ctx, uiduser, payload)

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
		Roles:     roles,
	}

	return result, sc, msg, nil
}

func (s *service) Logout(c echo.Context, ctx context.Context, uiduser string) (string, int, error) {
	// 1. Check Account
	users, sc, msg, err := s.UserRepository.MyAccount(ctx, uiduser)
	if err != nil {
		log.Print("Error "+msg, sc)
		return msg, sc, err
	}

	// 2. Delete Session
	msg, sc, err = s.SessionRepository.LogoutSession(c, ctx, users)
	if err != nil {
		log.Print("Error "+msg, sc)
		return msg, sc, err
	}

	// 3. delete fcm
	msgfcm, scfcm, err := s.FcmTokenRepository.LogoutFCMTokenUser(c, ctx, users.UidUser.String())
	if err != nil {
		log.Print("Error "+msgfcm, sc)
		return msgfcm, scfcm, err
	}

	return msg, sc, err
}
