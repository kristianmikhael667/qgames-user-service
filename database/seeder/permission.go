package seeder

import (
	"log"
	model "main/internal/model/users"
	"time"

	"gorm.io/gorm"
)

func permissionSeeder(db *gorm.DB) {
	now := time.Now()

	var permission = []model.Permission{
		{
			Name:      "Create Users",
			Slug:      "create-users",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Update Users",
			Slug:      "update-users",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "List Users",
			Slug:      "list-users",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Banned Users",
			Slug:      "banned-users",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Create Product",
			Slug:      "create-product",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Update Product",
			Slug:      "update-product",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "List Product",
			Slug:      "list-product",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Delete Product",
			Slug:      "delete-product",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Create Transaction",
			Slug:      "create-Transaction",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Update Transaction",
			Slug:      "update-Transaction",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Report Transaction",
			Slug:      "report-Transaction",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Create Payment",
			Slug:      "create-payment",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Topup Wallet",
			Slug:      "topup-wallet",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Check Wallet",
			Slug:      "check-wallet",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	if err := db.Create(&permission).Error; err != nil {
		log.Printf("Can't seeder data permission, with error %v \n", err)
	}
	log.Println("Success seed data permission")
}
