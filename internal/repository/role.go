package repository

import (
	"context"
	dto "main/internal/dto"
	model "main/internal/model"

	"gorm.io/gorm"
)

type Role interface {
	Save(ctx context.Context, roles *dto.RoleRequestBody) (model.Role, error)
	ExistByName(ctx context.Context, name *string) (bool, error)
}

type role struct {
	Db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *role {
	return &role{
		db,
	}
}

func (r *role) Save(ctx context.Context, roles *dto.RoleRequestBody) (model.Role, error) {
	newRoles := model.Role{
		Name:   roles.Name,
		Desc:   roles.Desc,
		Data:   roles.Data,
		Status: roles.Status,
	}

	if err := r.Db.WithContext(ctx).Save(&newRoles).Error; err != nil {
		return newRoles, nil
	}
	return newRoles, nil
}

func (r *role) ExistByName(ctx context.Context, name *string) (bool, error) {
	var (
		count   int64
		isExist bool
	)

	if err := r.Db.WithContext(ctx).Model(&model.Role{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}
