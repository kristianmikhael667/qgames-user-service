package repository

import (
	"context"
	"fmt"
	"main/helper"
	dto "main/internal/dto"
	model "main/internal/model"
	pkgdto "main/package/dto"
	pkgutil "main/package/util"
	"main/package/util/response"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.User, *pkgdto.PaginationInfo, error)
	Save(ctx context.Context, users *dto.RegisterUsersRequestBody) (model.User, error)
	ExistByEmail(ctx context.Context, email *string) (bool, error)
	ExistByPhone(ctx context.Context, email string) (bool, error)
	RequestOtp(ctx context.Context, phone string) (model.User, bool, string, error)
	VerifyOtp(ctx context.Context, phone string, otps string) (model.User, bool, string, error)
	GetAssignUsers(ctx context.Context, uidusers string) ([]model.Assign, error)
	UpdateAttemptOtp(ctx context.Context, phone string) (int16, string, error)
	UpdateAccount(ctx context.Context, uuid string, users *dto.UpdateUsersReqBody) (model.User, int16, string, error)
	LoginByPin(ctx context.Context, loginpin *dto.LoginByPin) (model.User, int16, string, error)
	LoginAdmin(ctx context.Context, loginadmin *dto.LoginAdmin) (model.User, int, string, error)
	MyAccount(ctx context.Context, iduser string) (model.User, int16, string, error)
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
			fmt.Println("error apa ini ", err.Error())
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

	if err := r.Db.WithContext(ctx).Model(&model.Otp{}).Where("phone = ?", phones).Order("created_at DESC").First(&otp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return users, false, "Sorry, your OTP has expired", err
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
		return users, false, "Expired Otp", nil
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

	return users, true, "Success verify OTP", nil
}

func (r *user) GetAssignUsers(ctx context.Context, uidusers string) ([]model.Assign, error) {
	var assign []model.Assign

	if err := r.Db.WithContext(ctx).Where("users = ? ", uidusers).Find(&assign).Error; err != nil {
		helper.Logger("error", "Assign Not Found", "Rc: "+string(rune(404)))
	}
	return assign, nil
}

func (r *user) UpdateAttemptOtp(ctx context.Context, phone string) (int16, string, error) {
	var attemp model.Attempt

	phones := strings.Replace(phone, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ?", phones).First(&attemp).Error; err != nil {
		return 404, "Error get phone attemp", err
	}
	attemp.OtpAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&attemp).Error; err != nil {
		return 403, "Error update phone attemp", err
	}
	return 201, "Success update", nil
}

func (r *user) UpdateAccount(ctx context.Context, uuid string, users *dto.UpdateUsersReqBody) (model.User, int16, string, error) {
	var user model.User

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("uid_user = ?", uuid).First(&user).Error; err != nil {
		return user, 403, "Error get users", err
	}

	// Check PIN
	if len(users.Pin) != 6 || !helper.IsNumeric(users.Pin) {
		return user, 403, "Invalid PIN", nil
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(users.Pin), bcrypt.DefaultCost)
	if err != nil {
		return user, 403, "Failed Generate Pin", err
	}

	user.Fullname = users.Fullname
	user.Email = users.Email
	user.Address = users.Address
	user.Pin = string(hashedPin)

	if err := r.Db.WithContext(ctx).Save(&user).Error; err != nil {
		return user, 400, "Failed Update User", err
	}

	return user, 201, "Success Update User", nil
}

func (r *user) LoginByPin(ctx context.Context, loginpin *dto.LoginByPin) (model.User, int16, string, error) {
	var user model.User
	var trylimit model.Attempt

	phones := strings.Replace(loginpin.Phone, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	// Validate Phone Number
	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).First(&user).Error; err != nil {
		helper.Logger("error", "Number Phone Not Found Users", "Rc: "+string(rune(404)))
		return user, 404, "Number Phone Not Found Users", err
	}

	// Validate TryLimit
	timenow := time.Now()
	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).First(&trylimit).Error; err != nil {
		// Create Limit
		newAttemp := model.Attempt{
			Phone:       phones,
			PinAttempt:  0,
			OtpAttempt:  0,
			LastAttempt: timenow,
		}
		if err := r.Db.WithContext(ctx).Save(&newAttemp).Error; err != nil {
			return user, 403, "Failed create attemp", err
		}
	}

	if (!helper.VerifyPin(loginpin.Pin, user.Pin)) && (trylimit.PinAttempt < 3) {
		helper.Logger("error", "Wrong PIN", "Rc: "+string(rune(403)))
		trylimit.PinAttempt = trylimit.PinAttempt + 1
		if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
			return user, 403, "Failed update trylimit", err
		}
		return user, 403, "Wrong PIN", response.CustomErrorBuilder(403, "Error", "Wrong PIN")

	} else if trylimit.PinAttempt == 3 {
		return user, 400, "Your access has been restricted, redirecting to PIN verification", response.CustomErrorBuilder(403, "Error", "Your access has been restricted, redirecting to PIN verification")
	}
	trylimit.PinAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return user, 403, "Failed update trylimit", err
	}

	return user, 201, "Success Login By Pin", nil
}

func (r *user) LoginAdmin(ctx context.Context, loginadmin *dto.LoginAdmin) (model.User, int, string, error) {
	var user model.User
	// Check email
	if err := r.Db.WithContext(ctx).Where("email = ?", loginadmin.Email).First(&user).Error; err != nil {
		helper.Logger("error", "Email Not Found", "Rc: "+string(rune(404)))
		return user, 404, "Email not found", err
	} else if !helper.VerifyPassword(loginadmin.Password, user.Password) {
		helper.Logger("error", "Password is wrong", "Rc: "+string(rune(404)))
		return user, 404, "Password not found", err
	} else {
		return user, 201, "Success Login Admin", nil
	}
}

func (r *user) MyAccount(ctx context.Context, iduser string) (model.User, int16, string, error) {
	var user model.User

	if err := r.Db.WithContext(ctx).Where("uid_user = ? ", iduser).Find(&user).Error; err != nil {
		helper.Logger("error", "Assign Not Found", "Rc: "+string(rune(404)))
		return user, 404, "User Not Found", nil

	}
	return user, 200, "Get User " + user.Fullname, nil
}
