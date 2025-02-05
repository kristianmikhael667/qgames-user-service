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

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.User, *pkgdto.PaginationInfo, error)
	FindIDUser(ctx context.Context, uid string) (model.User, error)
	Save(ctx context.Context, users *dto.RegisterUsersRequestBody) (model.User, error)
	ExistByEmail(ctx context.Context, email *string) (bool, error)
	ExistByPhone(ctx context.Context, email string) (bool, error)
	CreateUsers(ctx context.Context, phone string) (model.User, bool, int, string, error)
	CheckUser(ctx context.Context, reqOtp bool, phone string) (model.User, int, bool, string, error)
	VerifyOtp(ctx context.Context, phone string, otps string) (model.User, bool, string, error)
	UpdateAccount(ctx context.Context, users *dto.UpdateUsersReqBody) (model.User, int, string, error)
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

func (r *user) FindIDUser(ctx context.Context, uid string) (model.User, error) {
	var users model.User

	q := r.Db.WithContext(ctx).Model(&model.User{}).Where("uid_user = ?", uid).First(&users)

	err := q.Error

	return users, err
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
	// phones = strings.Replace(phones, "62", "0", -1)

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

func (r *user) CreateUsers(ctx context.Context, phone string) (model.User, bool, int, string, error) {
	var users model.User

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ? ", phone).First(&users).Error; err != nil {
		newUsers := model.User{
			Phone:  phone,
			Status: "active",
		}

		if err := r.Db.WithContext(ctx).Save(&newUsers).Error; err != nil {
			log.Print("Failed Create User With Number : "+phone, 400)
		} else {
			log.Print("Success Create User With Number: "+phone, 201)
			return newUsers, true, 205, "New User Successfully Created, Please Register", nil
		}
	}
	return users, false, 201, "Users valid with pin", nil
}

func (r *user) CheckUser(ctx context.Context, reqOtp bool, phone string) (model.User, int, bool, string, error) {
	var users model.User

	if reqOtp {
		if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ? ", phone).First(&users).Error; err != nil {
			return users, 201, true, "An instruction to verify your phone number has been sent to your phone.", nil
		}
	} else {
		if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ? ", phone).First(&users).Error; err != nil {
			return users, 404, false, "Users not found with checkuser", nil
		}
	}

	if users.Fullname == "" && users.Pin == "" && users.Email == "" && users.Address == "" {
		// user not complate regist, ketika input nomor lagi, maka akan diarahkan ke page regist code 205 Reset Content
		return users, 205, true, "Uncomplate register users, please full regist", nil
	} else {
		if reqOtp {
			return users, 200, false, "Users valid from request otp", nil
		}
		return users, 201, false, "Users valid with otp", nil
	}
}

func (r *user) VerifyOtp(ctx context.Context, phone string, otps string) (model.User, bool, string, error) {
	var otp model.Otp
	var users model.User

	phones := strings.Replace(phone, "+62", "0", -1)
	// phones = strings.Replace(phones, "62", "0", -1)

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
		log.Print("Expired OTP ", 403)
		return users, false, "Expired Otp", nil
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

	return users, true, "Success verify OTP", nil
}

func (r *user) UpdateAccount(ctx context.Context, users *dto.UpdateUsersReqBody) (model.User, int, string, error) {
	var user model.User

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", users.Phone).First(&user).Error; err != nil {
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
		return user, 400, "Failed Created User", err
	}

	return user, 201, "Success Created User", nil
}

func (r *user) LoginByPin(ctx context.Context, loginpin *dto.LoginByPin) (model.User, int, string, error) {
	var user model.User
	var trylimit model.Attempt

	phones := strings.Replace(loginpin.Phone, "+62", "0", -1)
	// phones = strings.Replace(phones, "62", "0", -1)

	// Validate Phone Number
	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).First(&user).Error; err != nil {
		log.Print("Number Phone Not Found Users ", 404)
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
		log.Print("Wrong PIN ", 403)
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
	// phones = strings.Replace(phones, "62", "0", -1)

	// Validate Phone Number
	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).First(&user).Error; err != nil {
		log.Print("Number Phone Not Found Users ", 404)
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
		log.Print("Wrong PIN ", 403)
		trylimit.PinAttempt = trylimit.PinAttempt + 1
		if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
			return false, 500, err
		}
		return false, 401, nil

	} else if trylimit.PinAttempt == 3 {
		log.Print("Error PIN 3 Times", 403)
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
		log.Print("Email Not Found", 404)
		return user, 404, "Email not found", err
	} else if !helper.VerifyPassword(loginadmin.Password, user.Password) {
		log.Print("Password is wrong", 404)
		return user, 404, "Password not found", err
	} else {
		return user, 201, "Success Login Admin", nil
	}
}

func (r *user) GetUserByNumber(ctx context.Context, phone string) (model.User, int, string, error) {
	var user model.User

	phones := strings.Replace(phone, "+62", "0", -1)
	// phones = strings.Replace(phones, "62", "0", -1)

	if err := r.Db.WithContext(ctx).Where("phone = ? ", phones).Find(&user).Error; err != nil {
		log.Print("User not found", 404)
		return user, 404, "User Not Found", err
	}
	return user, 200, "Get User " + user.Fullname, nil
}

func (r *user) MyAccount(ctx context.Context, iduser string) (model.User, int, string, error) {
	var user model.User

	if err := r.Db.WithContext(ctx).Where("uid_user = ? ", iduser).Find(&user).Error; err != nil {
		log.Print("User Not Found", 404)
		return user, 404, "User Not Found", nil

	}
	return user, 200, "Get User " + user.Fullname, nil
}

func (r *user) ResetPin(ctx context.Context, uid_user string, payload *dto.ConfirmPin) (model.User, int, string, error) {
	var users model.User
	if err := r.Db.WithContext(ctx).Where("uid_user = ? ", uid_user).First(&users).Error; err != nil {
		log.Print("User Not Found", 404)
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
	fmt.Println("lah kocal ", users.Phone)
	if err := r.Db.WithContext(ctx).Where("phone = ? ", users.Phone).First(&trylimit).Error; err != nil {
		return users, 404, "Not Found User in TryLimit", err
	}
	trylimit.PinAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return users, 500, "Failed update trylimit", err
	}
	return users, 201, "Success Update PIN User", nil
}
