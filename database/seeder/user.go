package seeder

import (
	"log"
	model "main/internal/model/users"
	"time"

	"gorm.io/gorm"
)

func userSeeder(db *gorm.DB) {
	now := time.Now()
	var users = []model.User{
		{
			Fullname: "Mikhael Developer",
			Email:    "mikhael@developer.com",
			Password: "alibi123",
			Common: model.Common{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			Fullname: "Mikhael Fullstack",
			Email:    "mikhael@fullstack.com",
			Password: "alibi123",
			Common: model.Common{
				ID:        2,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	if err := db.Create(&users).Error; err != nil {
		log.Printf("Cannot seeder data user, with error %v \n", err)
	}
	log.Println("Success seed data user")
}
