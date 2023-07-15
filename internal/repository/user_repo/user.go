package repository

import (
	"context"
	dto "main/internal/dto/users_req_res"
	model "main/internal/model/users"
	pkgdto "main/package/dto"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type User interface {
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.User, *pkgdto.PaginationInfo, error)
	Save(ctx context.Context, users *dto.RegisterUsersRequestBody) (model.User, error)
	ExistByEmail(ctx context.Context, email *string) (bool, error)
	ExistByPhone(ctx context.Context, email *string) (bool, error)
}

type user struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *user {
	return &user{
		db,
	}
}

func (r *user) FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.User, *pkgdto.PaginationInfo, error) {
	var users []model.User
	var count int64

	query := r.Db.WithContext(ctx).Model(&model.User{})

	if payload.Search != "" {
		search := "%" + strings.ToLower(payload.Search) + "%"
		query = query.Where("lower(fullname) LIKE ? or lower(email) Like ? ", search, search)
	}

	countQuery := query
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	limit, offset := pkgdto.GetLimitOffset(pagination)

	err := query.Limit(limit).Offset(offset).Find(&users).Error

	return users, pkgdto.CheckInfoPagination(pagination, count), err
}

func (r *user) Save(ctx context.Context, users *dto.RegisterUsersRequestBody) (model.User, error) {
	newUsers := model.User{
		Fullname: users.Fullname,
		Phone:    users.Phone,
		Email:    users.Email,
		Password: users.Password,
		Address:  users.Address,
		Profile:  users.Profile,
		Pin:      users.Pin,
	}
	if err := r.Db.WithContext(ctx).Save(&newUsers).Error; err != nil {
		return newUsers, nil
	}
	return newUsers, nil
}

func (r *user) ExistByEmail(ctx context.Context, email *string) (bool, error) {
	var (
		count   int64
		isExist bool
	)

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}

func (r *user) ExistByPhone(ctx context.Context, numbers *string) (bool, error) {
	phones := strings.Replace(*numbers, "62", "0", 2)
	re := regexp.MustCompile(`\+62(\D|$)`)
	result := re.ReplaceAllString(phones, "0")
	var (
		count   int64
		isExist bool
	)

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", result).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}
