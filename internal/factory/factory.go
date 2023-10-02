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
	FcmTokenRepository   repository.Fcmtoken
}

func NewFactory() *Factory {
	db := database.GetConnection()
	mongo := database.GetConnectionMongoDB()
	return &Factory{
		repository.NewUserRepository(db),
		repository.NewRoleRepository(db),
		repository.NewPermissionRepository(db),
		repository.NewAssign(db),
		repository.NewAttemptRepository(db),
		repository.NewOtpRepository(db),
		repository.NewSessionRepository(db),
		repository.NewFcmToken(mongo),
	}
}
