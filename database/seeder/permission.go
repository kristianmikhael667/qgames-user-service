package seeder

import (
	"log"
	model "main/internal/model"
	"time"

	"gorm.io/gorm"
)

func permissionSeeder(db *gorm.DB) {
	now := time.Now()

	var permission = []model.Permission{
		{
			Name:      "Common User",
			Slug:      "common-user",
			CreatedAt: now,
			UpdatedAt: now,
		},
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
			Name:      "List Product Default",
			Slug:      "list-product-default",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "List Product Basic",
			Slug:      "list-product-basic",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "List Product VIP",
			Slug:      "list-product-vip",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "List Product VVIP",
			Slug:      "list-product-vvip",
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
			Name:      "Create Trx",
			Slug:      "create-trx",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Update trx",
			Slug:      "update-trx",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:      "Report Trx",
			Slug:      "report-trx",
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
