package seeder

import (
	"log"
	model "main/internal/model"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func roleSeeder(db *gorm.DB) {
	now := time.Now()

	var roles = []model.Role{
		{
			UidRole: uuid.NewV4(),
			Name:    "Default",
			Desc:    "No Topup Balance",
			Data:    "user-default",
			Status:  "active",
			Common: model.Common{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			UidRole: uuid.NewV4(),
			Name:    "Basic",
			Desc:    "Min Top Up 300K",
			Data:    "user-basic",
			Status:  "active",
			Common: model.Common{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			UidRole: uuid.NewV4(),
			Name:    "VIP",
			Desc:    "Min Top Up 1.5Jt",
			Data:    "user-vip",
			Status:  "active",
			Common: model.Common{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			UidRole: uuid.NewV4(),
			Name:    "VVIP",
			Desc:    "min topup 3 jt ",
			Data:    "user-vvip",
			Status:  "active",
			Common: model.Common{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	if err := db.Create(&roles).Error; err != nil {
		log.Printf("Cannot seeder data role, with error %v \n", err)
	}
	log.Println("Success seed data role")
}
