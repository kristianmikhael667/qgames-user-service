package repository

import (
	"context"
	"main/helper"
	"main/internal/dto"
	model "main/internal/model"
	"main/internal/pkg/util"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Assign interface {
	FindUserID(ctx context.Context, users string) (model.Assign, error)
	Assign(ctx context.Context, users, role, permission string) error
	GetAssignUsers(ctx context.Context, uidusers string) ([]model.Assign, error)
	EditRolesTopup(c echo.Context, ctx context.Context, payload *dto.ReqAssign) (bool, error)
}

type assigns struct {
	Db *gorm.DB
}

func NewAssign(db *gorm.DB) *assigns {
	return &assigns{
		db,
	}
}

func (a *assigns) FindUserID(ctx context.Context, users string) (model.Assign, error) {
	var assign model.Assign

	q := a.Db.WithContext(ctx).Model(&model.Assign{}).Where("users = ?", users).First(&assign)

	err := q.Error

	return assign, err
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

func (r *assigns) EditRolesTopup(c echo.Context, ctx context.Context, payload *dto.ReqAssign) (bool, error) {
	var assign model.Assign
	authHeader := c.Request().Header.Get("Authorization")
	tokens, _ := util.ParseJWTToken(authHeader)

	if err := r.Db.WithContext(ctx).Where("users = ? ", tokens.Uuid).Find(&assign).Error; err != nil {
		helper.Logger("error", "Assign Users Not Found", "Rc: "+string(rune(404)))
		return false, err
	}

	if assign.Roles == "user-default" {
		if payload.PaymentFee >= 300000 {
			assign.Roles = "user-basic"
			if err := r.Db.WithContext(ctx).Save(&assign).Error; err != nil {
				return false, err
			}
		}
	} else if assign.Roles == "user-basic" {
		if payload.PaymentFee < 300000 {
			assign.Roles = "user-default"
			if err := r.Db.WithContext(ctx).Save(&assign).Error; err != nil {
				return false, err
			}
		} else if payload.PaymentFee >= 1500000 {
			assign.Roles = "user-vip"
			if err := r.Db.WithContext(ctx).Save(&assign).Error; err != nil {
				return false, err
			}
		}
	} else if assign.Roles == "user-vip" {
		if payload.PaymentFee < 1500000 {
			assign.Roles = "user-basic"
			if err := r.Db.WithContext(ctx).Save(&assign).Error; err != nil {
			} else if payload.PaymentFee >= 3000000 {
				assign.Roles = "user-vvip"
				return false, err
			}
		} else if assign.Roles == "user-vvip" {
			if payload.PaymentFee < 3000000 {
				assign.Roles = "user-vip"
			}
		}
	}
	return true, nil
}
