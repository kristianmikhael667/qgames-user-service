package dto

import "time"

type (
	UsersResponse struct {
		Uuid      string    `json:"uuid"`
		Fullname  string    `json:"fullname"`
		Phone     string    `json:"phone"`
		Email     string    `json:"email"`
		Address   string    `json:"address"`
		Profile   string    `json:"profile"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserResponseAll struct {
		Data  UsersResponse `json:"data"`
		Roles string        `json:"roles"`
	}

	UserWithJWTResponse struct {
		UsersResponse
		Token string `json:"token"`
		Admin bool   `json:"admin"`
	}

	UserWithCUDResponse struct {
		UsersResponse
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
