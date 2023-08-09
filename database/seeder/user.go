package seeder

import (
	"log"
	model "main/internal/model"
	pkgutil "main/package/util"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func userSeeder(db *gorm.DB) {
	now := time.Now()

	// Hash Password
	hashedPassword, _ := pkgutil.HashPassword("admin123")

	var users = []model.User{
		{
			UidUser:   uuid.NewV4(),
			Phone:     "081399941007",
			Fullname:  "Super Admin Qgames",
			Email:     "admin@gmail.com",
			Password:  hashedPassword,
			Pin:       hashedPassword,
			Address:   "QP Office, Grand Slipi Tower Lantai 5 Unit I. 1, Jalan S. Parman Kav 22-24 Jakarta Barat, DKI Jakarta",
			Profile:   "qgames.png",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
			Common: model.Common{
				ID:        1,
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
