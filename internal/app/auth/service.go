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
	utils "main/package/util"
	"main/package/util/response"
	"strings"

	"github.com/labstack/echo/v4"
)

type service struct {
	UserRepository     repository.User
	AssignRepository   repository.Assign
	AttemptRepository  repository.Attempt
	OtpRepository      repository.Otp
	SessionRepository  repository.Session
	FcmTokenRepository repository.Fcmtoken
}

type Service interface {
	RegisterUsers(ctx context.Context, payload *dto.RegisterUsersRequestBody) (*dto.UserWithJWTResponse, error)
	CheckPhone(ctx context.Context, payload *dto.RegisterUsersRequestBody) (bool, error)
	RequestOtp(c echo.Context, ctx context.Context, phone *dto.CheckPhoneReqBody) (string, int, bool, error)
	VerifyOtp(c echo.Context, ctx context.Context, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, string, int, error)
	LoginPin(c echo.Context, ctx context.Context, loginpin *dto.LoginByPin) (*dto.UserWithJWTResponse, string, int, error)
	LoginAdmin(ctx context.Context, loginadmin *dto.LoginAdmin) (*dto.UserWithJWTResponse, string, int, error)
	ConfirmReset(c echo.Context, ctx context.Context, phone *dto.CheckSession) (string, int, error)
	ResetDevice(c echo.Context, ctx context.Context, session *dto.ReqSessionReset) (*dto.UserWithJWTResponse, string, int, error)
	CheckPin(c echo.Context, ctx context.Context, token *dto.JWTClaims, loginpin *dto.CheckPin) (bool, int, string, error)
	RefreshToken(ctx context.Context, oldtoken *dto.JWTClaims) (*dto.UserWithJWTResponse, int, string, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository:     f.UserRepository,
		AssignRepository:   f.AssignRepository,
		AttemptRepository:  f.AttemptRepository,
		OtpRepository:      f.OtpRepository,
		SessionRepository:  f.SessionRepository,
		FcmTokenRepository: f.FcmTokenRepository,
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

func (s *service) RequestOtp(c echo.Context, ctx context.Context, phone *dto.CheckPhoneReqBody) (string, int, bool, error) {
	// Step 1. Number Check Regex
	phones := strings.Replace(phone.Phone, "+62", "0", -1)
	// phones = strings.Replace(phones, "62", "0", -1)

	// Step 2. Check Attempt
	trylimit, sc, msg, err := s.AttemptRepository.CreateAttempt(ctx, phones)

	if err != nil {
		return err.Error(), 500, false, err
	}

	// Step 3. Create Users
	users, sc, status, msg, err := s.UserRepository.CheckUser(ctx, true, phones)

	if err != nil {
		return err.Error(), sc, status, err
	}

	// Step 4. Create Session and Check Device Id
	msg, sc, otp, err := s.SessionRepository.CheckSession(c, ctx, users.UidUser.String(), phone.Phone, sc, msg)
	if err != nil {
		return err.Error(), sc, status, err
	}

	// Step 6. Create OTP and if send otp
	msg, sc, err = s.OtpRepository.SendOtp(ctx, phones, sc, otp, trylimit, msg)
	if err != nil {
		return err.Error(), sc, status, err
	}
	return msg, sc, status, nil
}

func (s *service) VerifyOtp(c echo.Context, ctx context.Context, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, string, int, error) {
	var result *dto.UserWithJWTResponse

	// If Number Tester
	if validotp.Phone == utils.Getenv("NUMBER_FAKE", "000") {
		datausers, sc, _, msg, err := s.UserRepository.CheckUser(ctx, false, utils.Getenv("NUMBER_FAKE", "000"))
		if err != nil {
			helper.Logger("error", msg+" User Tester", "Rc: "+string(rune(403)))
			return result, msg, sc, err
		}

		// Get all assign and loop
		response_assign, err := s.AssignRepository.GetAssignUsers(ctx, datausers.UidUser.String())

		if err != nil {
			helper.Logger("error", "Error get assign user service", "Rc: "+string(rune(403)))
		}

		helpers, sc, msg, err := helper.AuditOTPPlayStore(datausers, response_assign, result, validotp)
		if err != nil {
			return result, msg, sc, err
		}
		return helpers, msg, sc, nil
	}

	// Check OTP
	_, verifyOtp, msg, err := s.UserRepository.VerifyOtp(ctx, validotp.Phone, validotp.Otp)
	if err != nil {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, 403, err
	}
	if verifyOtp == false {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, 401, err
	}

	// Check number
	response, status, sc, msg, err := s.UserRepository.CreateUsers(ctx, validotp.Phone)

	// Create Assign Users only new user
	if status == true {
		err = s.AssignRepository.Assign(ctx, response.UidUser.String(), "user-default", "common-user,check-wallet,topup-wallet,list-product-default,create-trx")
		if err != nil {
			return result, msg, sc, err
		}
	}

	// Get all assign and loop
	response_assign, err := s.AssignRepository.GetAssignUsers(ctx, response.UidUser.String())

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
	statuscode, msg, err := s.AttemptRepository.UpdateAttemptOtp(ctx, validotp.Phone)
	if err != nil {
		return result, msg, statuscode, err
	}

	// Update Session
	msgSess, scode, errSess := s.SessionRepository.CreateSession(c, ctx, response.UidUser.String(), response.Phone, statuscode, msg)
	if errSess != nil {
		return result, msgSess, scode, err
	}

	// Step 5. Create FCM Token and Check FCM Token MongoDB
	msgFcm, err_fcm := s.FcmTokenRepository.CreateFCMTokenUser(c, ctx, response.UidUser.String())
	if err_fcm != nil {
		return result, msgFcm, 500, err_fcm
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
	return result, msg, sc, nil
}

func (s *service) LoginPin(c echo.Context, ctx context.Context, loginpin *dto.LoginByPin) (*dto.UserWithJWTResponse, string, int, error) {
	var result *dto.UserWithJWTResponse

	// Step 1. Check Number User
	users, sc, msg, err := s.UserRepository.GetUserByNumber(ctx, loginpin.Phone)
	if err != nil {
		return result, msg, sc, err
	}
	// Step 2. Check Session and Check Device Id
	msg, sc, _ = s.SessionRepository.CheckSessionPin(c, ctx, users.UidUser.String(), loginpin.Phone, sc, msg)
	if sc == 403 {
		return result, msg, sc, err
	}

	// Step 3. Login Pin
	responses, sc, msg, err := s.UserRepository.LoginByPin(ctx, loginpin)

	if err != nil {
		return result, msg, sc, err
	}

	// Step 4. Create FCM Token and Check FCM Token MongoDB
	_, err_fcm := s.FcmTokenRepository.CreateFCMTokenUser(c, ctx, users.UidUser.String())
	if err_fcm != nil {
		return result, "Error", 500, err_fcm
	}

	// Step 5. Get all assign and loop
	response_assign, err := s.AssignRepository.GetAssignUsers(ctx, responses.UidUser.String())

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
			Roles:     firstRole,
		},
		Token: token,
	}
	return result, msg, 201, nil
}

func (s *service) CheckPin(c echo.Context, ctx context.Context, token *dto.JWTClaims, loginpin *dto.CheckPin) (bool, int, string, error) {
	// Step 1. Check Number User
	users, scAcc, msgAcc, errAcc := s.UserRepository.MyAccount(ctx, token.Uuid)
	if errAcc != nil {
		return false, scAcc, msgAcc, errAcc
	}
	// Step 2. Check Session and Check Device Id
	msgSess, scSess, errSess := s.SessionRepository.CheckSessionPin(c, ctx, users.UidUser.String(), token.Phone, scAcc, msgAcc)
	fmt.Println("Aaa ", scSess)
	if scSess == 403 {
		return false, scSess, msgSess, errSess
	}

	// Step 3. Login Pin
	isPin, sc, err := s.UserRepository.CheckPin(ctx, token.Phone, loginpin.Pin)
	if sc == 401 {
		return isPin, sc, "Wrong PIN", err
	} else if sc == 403 {
		return isPin, sc, "Pin 3 Times", err
	} else {
		return isPin, sc, msgAcc, nil
	}
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

func (s *service) ConfirmReset(c echo.Context, ctx context.Context, phone *dto.CheckSession) (string, int, error) {
	// If Number Tester
	if phone.Phone == utils.Getenv("NUMBER_FAKE", "000") {
		_, sc, _, msg, err := s.UserRepository.CheckUser(ctx, false, utils.Getenv("NUMBER_FAKE", "000"))
		if err != nil {
			helper.Logger("error", msg+" User Tester", "Rc: "+string(rune(403)))
			return msg, int(sc), err
		}

		otp := utils.Getenv("OTP_FAKE", "")
		msg, scs, err := helper.AuditOTPDevicePlayStore(otp)
		if err != nil {
			return msg, scs, err
		}
		return msg, scs, nil
	}

	// Step 2. Check Attempt
	trylimit, scAttempt, msg, err := s.AttemptRepository.CreateAttempt(ctx, phone.Phone)
	if err != nil {
		return msg, scAttempt, err
	}

	// Step 3. Create Users
	users, scCheck, _, msg, err := s.UserRepository.CheckUser(ctx, true, phone.Phone)
	if err != nil {
		return msg, scCheck, err
	}

	// Step 3. Create Session and Check Device Id
	msg, scSess, errSess := s.SessionRepository.CheckSessionReset(c, ctx, users.UidUser.String(), phone)
	if errSess != nil {
		return msg, scSess, errSess
	}
	if scSess != 201 {
		return msg, scSess, errSess
	}

	otp := helper.GeneratePin(6)
	msg, scOtp, err := s.OtpRepository.SendOtp(ctx, phone.Phone, scSess, otp, trylimit, msg)
	if err != nil {
		return msg, scOtp, err
	}
	return "Send OTP Reset", scOtp, nil
}

func (s *service) ResetDevice(c echo.Context, ctx context.Context, session *dto.ReqSessionReset) (*dto.UserWithJWTResponse, string, int, error) {
	var result *dto.UserWithJWTResponse
	// If Number Tester
	if session.Phone == utils.Getenv("NUMBER_FAKE", "000") {
		datausers, sc, _, msg, err := s.UserRepository.CheckUser(ctx, false, utils.Getenv("NUMBER_FAKE", "000"))
		if err != nil {
			helper.Logger("error", msg+" User Tester", "Rc: "+string(rune(403)))
			return result, msg, int(sc), err
		}

		// Get all assign and loop
		response_assign, err := s.AssignRepository.GetAssignUsers(ctx, datausers.UidUser.String())
		if err != nil {
			helper.Logger("error", "Error get assign user service", "Rc: "+string(rune(403)))
		}

		// Step 1. Check Verify OTP
		helpers, sc, msg, err := helper.AuditResetDeviceOTP(datausers, response_assign, result, session)
		if err != nil {
			return helpers, msg, int(sc), err
		}

		// Step 2. Update Device ID
		msg, scs, err := s.SessionRepository.UpdateSession(c, ctx, datausers, session)
		if err != nil {
			helper.Logger("error", msg, "Rc: "+string(rune(403)))
			return result, msg, scs, err
		}
		return result, msg, scs, nil
	}

	// Step 1. Check Verify OTP
	_, verifyOtp, msg, err := s.UserRepository.VerifyOtp(ctx, session.Phone, session.Otp)
	if err != nil {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, 403, err
	}
	if verifyOtp == false {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msg, 403, err
	}

	// Get User
	responses, sc, msgUsr, err := s.UserRepository.GetUserByNumber(ctx, session.Phone)
	if err != nil {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msgUsr, sc, err
	}

	// Get all assign and loop
	response_assign, err := s.AssignRepository.GetAssignUsers(ctx, responses.UidUser.String())

	if err != nil {
		helper.Logger("error", "Error get assign user service", "Rc: "+string(rune(403)))
	}
	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(responses.UidUser.String(), responses.Email, responses.Phone, firstRole, permissions, false)

	// Update Limit
	statuscode, msg, err := s.AttemptRepository.UpdateAttemptOtp(ctx, responses.Phone)
	if statuscode != 201 && statuscode != 205 {
		return result, msg, int(statuscode), err
	}

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, msg, 401, err
	}

	// Step 2. Update Device ID
	msgSess, scSess, err := s.SessionRepository.UpdateSession(c, ctx, responses, session)
	if err != nil {
		helper.Logger("error", msg, "Rc: "+string(rune(403)))
		return result, msgSess, scSess, err
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

	return result, msg, scSess, nil
}

func (s *service) RefreshToken(ctx context.Context, oldtoken *dto.JWTClaims) (*dto.UserWithJWTResponse, int, string, error) {
	var result *dto.UserWithJWTResponse

	// Get Users
	users, sc, msg, err := s.UserRepository.MyAccount(ctx, oldtoken.Uuid)
	if err != nil {
		return result, sc, msg, err
	}

	response_assign, err := s.AssignRepository.GetAssignUsers(ctx, users.UidUser.String())

	if err != nil {
		helper.Logger("error", "Error get assign user service", "Rc: "+string(rune(403)))
	}

	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(users.UidUser.String(), users.Email, users.Phone, firstRole, permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, sc, msg, err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      users.UidUser.String(),
			Fullname:  users.Fullname,
			Phone:     users.Phone,
			Email:     users.Email,
			Address:   users.Address,
			Profile:   users.Profile,
			CreatedAt: users.CreatedAt,
			UpdatedAt: users.UpdatedAt,
			Roles:     firstRole,
		},
		Token: token,
	}
	return result, 201, msg, nil
}
