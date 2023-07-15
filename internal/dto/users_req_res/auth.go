package dto

import (
	"github.com/golang-jwt/jwt/v4"
)

type (
	RegisterUsersRequestBody struct {
		Fullname string `json:"fullname" validate:"omitempty"`
		Phone    string `json:"phone" validate:"omitempty"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"omitempty"`
		Pin      string `json:"pin" validate:"omitempty"`
		Address  string `json:"address" validate:"omitempty"`
		Profile  string `json:"profile"`
	}

	CheckPhoneReqBody struct {
		Phone string `json:"phone" validate:"omitempty"`
	}

	JWTClaims struct {
		Uuid  string `json:"uid_user"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		jwt.RegisteredClaims
	}
)
