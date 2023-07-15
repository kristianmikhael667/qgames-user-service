package seeder

import (
	"main/database"

	"gorm.io/gorm"
)

type seed struct {
	DB *gorm.DB
}

func NewSeeder() *seed {
	return &seed{database.GetConnection()}
}

func (s *seed) SeedAll() {
	// userSeeder(s.DB)
	roleSeeder(s.DB)
}

func (s *seed) DeleteAll() {
	s.DB.Exec("DELETE FROM user")
}
