package repository

import (
	"context"
	"main/helper"
	model "main/internal/model"

	"gorm.io/gorm"
)

type Assign interface {
	Assign(ctx context.Context, users, role, permission string) error
	GetAssignUsers(ctx context.Context, uidusers string) ([]model.Assign, error)
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

func (r *assigns) GetAssignUsers(ctx context.Context, uidusers string) ([]model.Assign, error) {
	var assign []model.Assign

	if err := r.Db.WithContext(ctx).Where("users = ? ", uidusers).Find(&assign).Error; err != nil {
		helper.Logger("error", "Assign Not Found", "Rc: "+string(rune(404)))
	}
	return assign, nil
}
