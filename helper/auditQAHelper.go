package helper

import (
	"main/internal/dto"
	"main/internal/model"
	"main/internal/pkg/util"
	utils "main/package/util"
)

func AuditOTPPlayStore(datausers model.User, response_assign []model.Assign, result *dto.UserWithJWTResponse, validotp *dto.RequestPhoneOtp) (*dto.UserWithJWTResponse, int, string, error) {
	if validotp.Otp != utils.Getenv("OTP_FAKE", "000000") {
		return result, 401, "Wrong OTP Tester", nil
	}

	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(datausers.UidUser.String(), datausers.Email, datausers.Phone, firstRole, permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, 401, "Error Create JWT Token", err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      datausers.UidUser.String(),
			Fullname:  datausers.Fullname,
			Phone:     datausers.Phone,
			Email:     datausers.Email,
			Address:   datausers.Address,
			Profile:   datausers.Profile,
			CreatedAt: datausers.CreatedAt,
			UpdatedAt: datausers.UpdatedAt,
		},
		Token: token,
	}
	return result, 201, "Success OTP Helper Audit", nil
}

func AuditOTPDevicePlayStore(otp string) (string, int, error) {
	if otp != utils.Getenv("OTP_FAKE", "000000") {
		return "Wrong OTP Tester", 401, nil
	}
	return "Send OTP Reset", 201, nil
}

func AuditResetDeviceOTP(datausers model.User, response_assign []model.Assign, result *dto.UserWithJWTResponse, validotp *dto.ReqSessionReset) (*dto.UserWithJWTResponse, int, string, error) {
	if validotp.Otp != utils.Getenv("OTP_FAKE", "000000") {
		return result, 401, "Wrong OTP Tester", nil
	}

	firstRole := response_assign[0].Roles

	var permissions []string
	for _, assign := range response_assign {
		permissions = append(permissions, assign.Permissions)
	}

	claims := util.CreateJWTClaims(datausers.UidUser.String(), datausers.Email, datausers.Phone, firstRole, permissions, false)

	token, err := util.CreateJWTToken(claims)
	if err != nil {
		return result, 401, "Error Create JWT Token", err
	}

	result = &dto.UserWithJWTResponse{
		UsersResponse: dto.UsersResponse{
			Uuid:      datausers.UidUser.String(),
			Fullname:  datausers.Fullname,
			Phone:     datausers.Phone,
			Email:     datausers.Email,
			Address:   datausers.Address,
			Profile:   datausers.Profile,
			CreatedAt: datausers.CreatedAt,
			UpdatedAt: datausers.UpdatedAt,
		},
		Token: token,
	}
	return result, 201, "Success OTP Helper Audit", nil
}
