package repository

import (
	"context"
	"fmt"
	"main/helper"
	dto "main/internal/dto/users_req_res"
	model "main/internal/model/users"
	pkgdto "main/package/dto"
	pkgutil "main/package/util"
	"main/package/util/response"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User interface {
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.User, *pkgdto.PaginationInfo, error)
	Save(ctx context.Context, users *dto.RegisterUsersRequestBody) (model.User, error)
	ExistByEmail(ctx context.Context, email *string) (bool, error)
	ExistByPhone(ctx context.Context, email string) (bool, error)
	RequestOtp(ctx context.Context, phone string) (model.User, bool, string, error)
	VerifyOtp(ctx context.Context, phone string, otps string) (model.User, bool, string, error)
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
		return newUsers, err
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

func (r *user) ExistByPhone(ctx context.Context, numbers string) (bool, error) {
	phones := strings.Replace(numbers, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	var (
		count   int64
		isExist bool
	)

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", phones).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}

func (r *user) RequestOtp(ctx context.Context, phone string) (model.User, bool, string, error) {
	var users model.User
	var trylimit model.Attempt

	phones := strings.Replace(phone, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	timenow := time.Now()
	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ? ", phones).First(&trylimit).Error; err != nil {
		// create try limit
		newAttemp := model.Attempt{
			Phone:       phones,
			PinAttempt:  0,
			OtpAttempt:  0,
			LastAttempt: timenow,
		}
		if err := r.Db.WithContext(ctx).Save(&newAttemp).Error; err != nil {
			return users, false, "Failed create attemp", err
		}
	}

	curr := timenow.Format("2006-01-02")
	lastTest := trylimit.LastAttempt.Format("2006-01-02")

	if (trylimit.OtpAttempt >= 3) && (curr == lastTest) {
		helper.Logger("error", "Otp already 3 times  : "+string(rune(403)), "Rc: "+string(rune(403)))
		return users, false, "Otp already 3 times", response.CustomErrorBuilder(403, "Error", "Otp already 3 times")

	} else if curr != lastTest {
		trylimit.OtpAttempt = 0
	}

	otp := helper.GeneratePin(6)
	status_user := false

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ? ", phones).First(&users).Error; err != nil {
		status_user = true

		newUsers := model.User{
			Phone: phones,
		}

		if err := r.Db.WithContext(ctx).Save(&newUsers).Error; err != nil {
			helper.Logger("error", "Failed Create User With Number : "+phones, "400")
		} else {
			helper.Logger("info", "Success Create User With Number: "+phones, "Rc: "+string(rune(201)))
			users.UidUser = newUsers.UidUser
		}
	}

	// Sementara karena blm ada otp dari vendor maka kita buat send ke logs dlu
	helper.Logger("info", "Success Send Otp with Number: "+phones+" Your OTP : "+otp, "Rc: "+string(rune(201)))

	hashedOtp, err := pkgutil.HashPassword(otp)
	if err != nil {
		return users, status_user, "Error Hash OTP", response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	newOtp := model.Otp{
		Phone:     phones,
		Otp:       hashedOtp,
		ExpiredAt: time.Now().Add(1 * time.Minute),
	}
	if err := r.Db.WithContext(ctx).Save(&newOtp).Error; err != nil {
		return users, false, "Failed create otp", err
	}

	if status_user == false {
		trylimit.OtpAttempt++
		trylimit.LastAttempt = time.Now()

		if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
			return users, false, "Failed create otp", err
		}
	}

	// Call API
	msg, boolean := helper.SendOtp(phones, otp)
	if boolean == false {
		helper.Logger("error", "Failed OTP : "+phones, "400")
	}
	helper.Logger("info", msg, "Rc: "+string(rune(201)))

	return users, status_user, "Success create/get users", nil
}

func (r *user) VerifyOtp(ctx context.Context, phone string, otps string) (model.User, bool, string, error) {
	var otp model.Otp
	var users model.User

	phones := strings.Replace(phone, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	var (
		isExist bool
	)

	if err := r.Db.WithContext(ctx).Model(&model.Otp{}).Where("phone = ?", phones).Order("created_at DESC").First(&otp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return users, isExist, "Phone number is not registered", response.CustomErrorBuilder(403, "Error", "Expired Otp")
		}
	}

	// Set 1 minute
	expiredminute := 1

	currentTime := time.Now()
	expiredOtp := otp.ExpiredAt

	diff := expiredOtp.Sub(currentTime)
	minutesPassed := int(diff.Seconds())

	if minutesPassed <= expiredminute {
		helper.Logger("error", "Expired Otp : "+string(rune(403)), "Rc: "+string(rune(403)))
		return users, isExist, "Expired Otp", response.CustomErrorBuilder(403, "Error", "Expired Otp")
	}

	// Get users
	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ? ", phones).First(&users).Error; err != nil {
		helper.Logger("error", "Number not found", "Rc: "+string(rune(403)))
		return users, false, "Number not found", err
	}
	// Compare OTP
	check := helper.VerifyOtp(phones, otps, otp.Otp)
	if check == false {
		return model.User{}, false, "Failed verify otp", nil
	}

	// If success delete otp
	result := r.Db.Debug().Unscoped().Where("phone = ?", phones).Delete(&model.Otp{})
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	}
	// Update User
	users.Status = "active"
	if err := r.Db.WithContext(ctx).Save(&users).Error; err != nil {
		return users, false, "Failed update status user", err
	}

	return users, isExist, "Success verify OTP", nil
}
