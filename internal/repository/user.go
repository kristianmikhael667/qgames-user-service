package repository

import (
	"context"
	"fmt"
	"main/helper"
	dto "main/internal/dto"
	model "main/internal/model"
	pkgdto "main/package/dto"
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
	CreateUsers(ctx context.Context, phone string, device_id string) (model.User, int, bool, string, error)
	CheckUser(ctx context.Context, phone string) (model.User, int16, bool, string, error)
	VerifyOtp(ctx context.Context, phone string, otps string) (model.User, bool, string, error)
	UpdateAccount(ctx context.Context, uuid string, users *dto.UpdateUsersReqBody) (model.User, int16, string, error)
	LoginByPin(ctx context.Context, loginpin *dto.LoginByPin) (model.User, int, string, error)
	CheckPin(ctx context.Context, phone string, loginpin string) (bool, int, error)
	LoginAdmin(ctx context.Context, loginadmin *dto.LoginAdmin) (model.User, int, string, error)
	GetUserByNumber(ctx context.Context, phone string) (model.User, int, string, error)
	MyAccount(ctx context.Context, iduser string) (model.User, int, string, error)
	ResetPin(ctx context.Context, uid_user string, payload *dto.ConfirmPin) (model.User, int, string, error)
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

func (r *user) CreateUsers(ctx context.Context, phone string, device_id string) (model.User, int, bool, string, error) {
	var users model.User

	status_user := false

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ? ", phone).First(&users).Error; err != nil {
		status_user = true

		newUsers := model.User{
			Phone: phone,
		}

		if err := r.Db.WithContext(ctx).Save(&newUsers).Error; err != nil {
			fmt.Println("error apa ini ", err.Error())
			helper.Logger("error", "Failed Create User With Number : "+phone, "400")
		} else {
			helper.Logger("info", "Success Create User With Number: "+phone, "Rc: "+string(rune(201)))
			users.UidUser = newUsers.UidUser
		}
	}

	fmt.Println("ssss ", users.Status)
	if users.Status != "active" {
		return users, 201, status_user, "An instruction to verify your phone number has been sent to your phone.", nil
	} else if users.Fullname == "" && users.Pin == "" && users.Email == "" && users.Address == "" {
		// user not complate regist, ketika input nomor lagi, maka akan diarahkan ke page regist code 205 Reset Content
		return users, 205, status_user, "Uncomplate register users, please full regist", nil
	} else {
		return users, 200, false, "Users valid with pin", nil
	}
}

func (r *user) CheckUser(ctx context.Context, phone string) (model.User, int16, bool, string, error) {
	var users model.User
	if users.Fullname == "" && users.Pin == "" && users.Email == "" && users.Address == "" {
		// user not complate regist, ketika input nomor lagi, maka akan diarahkan ke page regist code 205 Reset Content
		return users, 205, true, "Uncomplate register users, please full regist", nil
	} else {
		return users, 201, false, "Users valid with pin", nil
	}
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

	// Set 2 minute
	expiredminute := 2

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

func (r *user) LoginByPin(ctx context.Context, loginpin *dto.LoginByPin) (model.User, int, string, error) {
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
			return user, 500, "Failed create attemp", err
		}
	}

	if (!helper.VerifyPin(loginpin.Pin, user.Pin)) && (trylimit.PinAttempt < 3) {
		helper.Logger("error", "Wrong PIN", "Rc: "+string(rune(403)))
		trylimit.PinAttempt = trylimit.PinAttempt + 1
		if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
			return user, 500, "Failed update trylimit", err
		}
		return user, 401, "Wrong PIN", response.CustomErrorBuilder(401, "Error", "Wrong PIN")

	} else if trylimit.PinAttempt == 3 {
		return user, 400, "Your access has been restricted, redirecting to PIN verification", response.CustomErrorBuilder(400, "Error", "Your access has been restricted, redirecting to PIN verification")
	}
	trylimit.PinAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return user, 500, "Failed update trylimit", err
	}

	return user, 201, "Success Login By Pin", nil
}

func (r *user) CheckPin(ctx context.Context, phone string, loginpin string) (bool, int, error) {
	var user model.User
	var trylimit model.Attempt

	phones := strings.Replace(phone, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	// Validate Phone Number
	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).First(&user).Error; err != nil {
		helper.Logger("error", "Number Phone Not Found Users", "Rc: "+string(rune(404)))
		return false, 404, err
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
			return false, 500, err
		}
	}

	if (!helper.VerifyPin(loginpin, user.Pin)) && (trylimit.PinAttempt < 3) {
		helper.Logger("error", "Wrong PIN", "Rc: "+string(rune(403)))
		trylimit.PinAttempt = trylimit.PinAttempt + 1
		if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
			return false, 500, err
		}
		return false, 401, nil

	} else if trylimit.PinAttempt == 3 {
		helper.Logger("error", "Pin 3 times", "Rc: "+string(rune(403)))
		return false, 403, nil
	}
	trylimit.PinAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return false, 500, err
	}

	return true, 201, nil
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

func (r *user) GetUserByNumber(ctx context.Context, phone string) (model.User, int, string, error) {
	var user model.User

	phones := strings.Replace(phone, "+62", "0", -1)
	phones = strings.Replace(phones, "62", "0", -1)

	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).Find(&user).Error; err != nil {
		helper.Logger("error", "User Not Found", "Rc: "+string(rune(404)))
		return user, 404, "User Not Found", err
	}
	return user, 200, "Get User " + user.Fullname, nil
}

func (r *user) MyAccount(ctx context.Context, iduser string) (model.User, int, string, error) {
	var user model.User

	if err := r.Db.WithContext(ctx).Where("uid_user = ? ", iduser).Find(&user).Error; err != nil {
		helper.Logger("error", "User Not Found", "Rc: "+string(rune(404)))
		return user, 404, "User Not Found", nil

	}
	return user, 200, "Get User " + user.Fullname, nil
}

func (r *user) ResetPin(ctx context.Context, uid_user string, payload *dto.ConfirmPin) (model.User, int, string, error) {
	var users model.User

	if err := r.Db.WithContext(ctx).Where("uid_user = ? ", uid_user).Find(&users).Error; err != nil {
		helper.Logger("error", "User Not Found", "Rc: "+string(rune(404)))
		return users, 404, "User Not Found", nil
	}

	// Change PIN
	if payload.NewPin != payload.ConfirmPin {
		return users, 401, "Different New Pin with Confirm New Pin", nil
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(payload.NewPin), bcrypt.DefaultCost)
	if err != nil {
		return users, 403, "Failed Generate Pin", err
	}

	users.Pin = string(hashedPin)

	if err := r.Db.WithContext(ctx).Save(&users).Error; err != nil {
		return users, 500, "Failed Update User", err
	}

	// Set Limit 0
	var trylimit model.Attempt
	if err := r.Db.WithContext(ctx).Where("phone = ? ", users.Phone).First(&trylimit).Error; err != nil {
		return users, 404, "Not Found User in TryLimit", err
	}
	trylimit.PinAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return users, 500, "Failed update trylimit", err
	}
	return users, 201, "Success Update PIN User", nil
}
