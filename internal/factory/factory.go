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
	AttemptRepository    repository.Attempt
	OtpRepository        repository.Otp
	SessionRepository    repository.Session
}

func NewFactory() *Factory {
	db := database.GetConnection()
	return &Factory{
		repository.NewUserRepository(db),
		repository.NewRoleRepository(db),
		repository.NewPermissionRepository(db),
		repository.NewAssign(db),
		repository.NewAttemptRepository(db),
		repository.NewOtpRepository(db),
		repository.NewSessionRepository(db),
	}
}
