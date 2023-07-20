package factory

import (
	"main/database"
	repository "main/internal/repository"
)

type Factory struct {
	UserRepository       repository.User
	RoleRepository       repository.Role
	PermissionRepository repository.Permission
	AssignRepository     repository.Assign
}

func NewFactory() *Factory {
	db := database.GetConnection()
	return &Factory{
		repository.NewUserRepository(db),
		repository.NewRoleRepository(db),
		repository.NewPermissionRepository(db),
		repository.NewAssign(db),
	}
}
