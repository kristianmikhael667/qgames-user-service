package factory

import (
	"main/database"
	"main/internal/repository"
)

type Factory struct {
	UserRepository repository.User
}

func NewFactory() *Factory {
	db := database.GetConnection()
	return &Factory{
		repository.NewUserRepository(db),
	}
}
