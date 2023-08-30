package helper

import (
	"main/internal/dto"
	"main/internal/pkg/util"
	utils "main/package/util"
	"strings"
	"time"
)

func AuditOTPPlayStore(result *dto.UserWithJWTResponse, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, int16, string, error) {
	if validotp.Otp != utils.Getenv("OTP_FAKE", "000000") {
		return result, 401, "Wrong OTP Tester", nil
	}
	response_assign := "common-user,check-wallet,topup-wallet,list-product-default,create-trx"

	permissions := strings.Split(response_assign, ",")

	for i, assign := range permissions {
		permissions[i] = strings.TrimSpace(assign)
	}

	claims := util.CreateJWTClaims(utils.Getenv("UUID_USER_FAKE", "000"), utils.Getenv("EMAIL_FAKE", "000"), utils.Getenv("NUMBER_FAKE", "000"), "user-default", permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, 401, "Error Create JWT Token", err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      utils.Getenv("UUID_USER_FAKE", "000"),
			Fullname:  utils.Getenv("FULLNAME_FAKE", "000"),
			Phone:     utils.Getenv("NUMBER_FAKE", "000"),
			Email:     utils.Getenv("EMAIL_FAKE", "000"),
			Address:   "JAKARTA",
			Profile:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Token: token,
	}
	return result, 201, "Success OTP Helper Audit", nil
}

func AuditPINPlayStore(result *dto.UserWithJWTResponse, validpin *dto.LoginByPin) (*dto.UserWithJWTResponse, int, string, error) {
	if validpin.Pin != utils.Getenv("OTP_FAKE", "000000") {
		return result, 401, "Wrong PIN Tester", nil
	}
	response_assign := "common-user,check-wallet,topup-wallet,list-product-default,create-trx"

	permissions := strings.Split(response_assign, ",")

	for i, assign := range permissions {
		permissions[i] = strings.TrimSpace(assign)
	}

	claims := util.CreateJWTClaims(utils.Getenv("UUID_USER_FAKE", "000"), utils.Getenv("EMAIL_FAKE", "000"), utils.Getenv("NUMBER_FAKE", "000"), "user-default", permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, 401, "Error Create JWT Token", err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      utils.Getenv("UUID_USER_FAKE", "000"),
			Fullname:  utils.Getenv("FULLNAME_FAKE", "000"),
			Phone:     utils.Getenv("NUMBER_FAKE", "000"),
			Email:     utils.Getenv("EMAIL_FAKE", "000"),
			Address:   "JAKARTA",
			Profile:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Token: token,
	}
	return result, 201, "Success PIN Helper Audit", nil
}

func AuditProfilePlayStore(result *dto.UsersResponse) (*dto.UsersResponse, int, string, error) {
	user_data := &dto.UsersResponse{
		Uuid:      utils.Getenv("UUID_USER_FAKE", "000"),
		Fullname:  utils.Getenv("FULLNAME_FAKE", "000"),
		Phone:     utils.Getenv("NUMBER_FAKE", "000"),
		Email:     utils.Getenv("EMAIL_FAKE", "000"),
		Address:   "JAKARTA",
		Profile:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Roles:     "user-default",
	}
	return user_data, 200, "Get User Helper Audit", nil
}
