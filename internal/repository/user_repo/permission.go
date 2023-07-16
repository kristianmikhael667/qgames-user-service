package repository

import (
	"context"
	dto "main/internal/dto/users_req_res"
	model "main/internal/model/users"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Permission interface {
	Save(ctx context.Context, permissions *dto.PermissionRequestBody) (model.Permission, error)
	ExistByNamePermission(ctx context.Context, name *string) (bool, error)
}

type permission struct {
	Db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *permission {
	return &permission{
		db,
	}
}

func (r *permission) Save(ctx context.Context, permissions *dto.PermissionRequestBody) (model.Permission, error) {
	slug := slug.Make(permissions.Name)
	newPermission := model.Permission{
		Name: permissions.Name,
		Slug: slug,
	}

	if err := r.Db.WithContext(ctx).Save(&newPermission).Error; err != nil {
		return newPermission, nil
	}
	return newPermission, nil
}

func (r *permission) ExistByNamePermission(ctx context.Context, name *string) (bool, error) {
	var (
		count   int64
		isExist bool
	)

	if err := r.Db.WithContext(ctx).Model(&model.Permission{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}
