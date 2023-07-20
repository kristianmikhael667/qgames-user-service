package repository

import (
	"context"
	model "main/internal/model"

	"gorm.io/gorm"
)

type Assign interface {
	Assign(ctx context.Context, users, role, permission string) error
}

type assigns struct {
	Db *gorm.DB
}

func NewAssign(db *gorm.DB) *assigns {
	return &assigns{
		db,
	}
}

func (a *assigns) Assign(ctx context.Context, users, role, permission string) error {
	assign := model.Assign{
		Users: users, Roles: role, Permissions: permission, Status: "active",
	}

	if err := a.Db.WithContext(ctx).Create(&assign).Error; err != nil {
		return err
	}
	return nil
}
