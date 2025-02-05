package dto

import (
	"github.com/golang-jwt/jwt/v4"
)

type (
	RegisterUsersRequestBody struct {
		Fullname string `json:"fullname" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Pin      string `json:"pin" validate:"required"`
		Address  string `json:"address" validate:"required"`
		Profile  string `json:"profile"`
	}

	UpdateUsersReqBody struct {
		Phone    string `json:"phone" validate:"required"`
		Fullname string `json:"fullname" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Address  string `json:"address" validate:"required"`
		Pin      string `json:"pin" validate:"required"`
	}

	LoginByPin struct {
		Phone string `json:"phone" validate:"required"`
		Pin   string `json:"pin" validate:"required"`
	}

	CheckPin struct {
		Pin string `json:"pin" validate:"required"`
	}

	RefreshToken struct {
		NewToken string `json:"new_token" validate:"required"`
	}

	LoginAdmin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	CheckPhoneReqBody struct {
		Phone string `json:"phone" validate:"required"`
	}

	PhoneAuditTester struct {
		Phone string `json:"phone" validate:"required"`
	}

	CheckSession struct {
		Phone string `json:"phone" validate:"required"`
		Reset string `json:"reset" validate:"required"`
	}

	ReqSessionReset struct {
		Phone string `json:"phone" validate:"required"`
		Otp   string `json:"otp" validate:"required"`
	}

	RequestPhoneOtp struct {
		Phone string `json:"phone" validate:"required"`
		Otp   string `json:"otp" validate:"required"`
	}

	JWTClaims struct {
		Uuid        string   `json:"uid_user"`
		Email       string   `json:"email"`
		Phone       string   `json:"phone"`
		Roles       string   `json:"roles"`
		Permissions []string `json:"permissions"`
		Admin       bool     `json:"admin"`
		jwt.RegisteredClaims
	}

	ConfirmPin struct {
		NewPin     string `json:"new_pin" validate:"required"`
		ConfirmPin string `json:"confirm_pin" validate:"required"`
	}

	RequestReset struct {
		Phone string `json:"phone" validate:"required"`
	}
)
