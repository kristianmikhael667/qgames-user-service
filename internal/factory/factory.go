package factory

import (
	"main/database"
	repository "main/internal/repository/user_repo"
)

type Factory struct {
	UserRepository repository.User
	RoleRepository repository.Role
}

func NewFactory() *Factory {
	db := database.GetConnection()
	return &Factory{
		repository.NewUserRepository(db),
		repository.NewRoleRepository(db),
	}
}
