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

	UserWithJWTResponse struct {
		UsersResponse
		Token string `json:"token"`
	}

	UserWithCUDResponse struct {
		UsersResponse
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
