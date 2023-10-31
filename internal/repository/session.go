package repository

import (
	"context"
	"main/helper"
	dto "main/internal/dto"
	model "main/internal/model"
	"main/package/util"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Session interface {
	CreateSession(ctx context.Context, uid_users string, device_id string, phone string, status int16, msg string) (string, int16, error)
	UpdateSession(ctx context.Context, sc int, msg string, session *dto.ReqSessionReset) (string, int, error)
	LogoutSession(ctx context.Context, phone string, device *dto.DeviceId) (string, int, error)
	CheckSession(ctx context.Context, uid_users string, device_id string, phone string, status int, msg string) (string, int, string, error)
}

type session struct {
	Db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *session {
	return &session{
		db,
	}
}

func (r *session) CreateSession(ctx context.Context, uid_users string, device_id string, phone string, status int16, msg string) (string, int16, error) {
	var sessions model.Session

	totalDevice := util.Getenv("TOTAL_DEVICE", "")
	intDevice, _ := strconv.ParseInt(totalDevice, 10, 16)
	int16Value := int16(intDevice)

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		// Create sesssion for new user initial
		newSession := model.Session{
			UserId:       uid_users,
			DeviceId:     device_id,
			LoginInAt:    time.Now(),
			TotalDevice:  1,
			Status:       true,
			ChangeDevice: nil,
			LoggedOutAt:  nil,
		}
		if err := r.Db.WithContext(ctx).Save(&newSession).Error; err != nil {
			return "Failed create session", 500, err
		}
		return msg, 201, nil
	}

	// Check all device id
	devices := strings.Split(sessions.DeviceId, ",")
	var isDevice bool
	for _, d := range devices {
		if d == device_id {
			isDevice = true
			break
		}
	}

	if isDevice == false && sessions.Status == true && sessions.LoggedOutAt == nil && status == 200 {
		// User device A sudah login, tetapi ada device B ingin login maka masih bisa
		sessions.TotalDevice = sessions.TotalDevice + 1
		if !strings.Contains(sessions.DeviceId, device_id) {
			sessions.DeviceId = sessions.DeviceId + "," + device_id
		}
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed create session", 500, err
		}
		return msg, status, nil
	} else if sessions.LoggedOutAt != nil && sessions.TotalDevice <= int16Value && sessions.Status == false {
		// User sudah logout di device a tetapi ingin login di device b
		sessions.Status = true
		sessions.DeviceId = device_id
		sessions.LoggedOutAt = nil
		sessions.LoginInAt = time.Now()
		sessions.TotalDevice = sessions.TotalDevice + 1
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed update session", 500, nil
		}
		return "Login OTP", 201, nil
	} else {
		// User logout di device yg sama dan user login dengan device yang sama
		sessions.Status = true
		sessions.LoggedOutAt = nil
		sessions.LoginInAt = time.Now()
		sessions.TotalDevice = sessions.TotalDevice + 1
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed update session", 500, nil
		}
		return msg, status, nil
	}

}

func (r *session) UpdateSession(ctx context.Context, sc int, msg string, session *dto.ReqSessionReset) (string, int, error) {
	var sessions model.Session
	var trylimit model.Attempt
	var users model.User
	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", session.Phone).First(&users).Error; err != nil {
		return "Phone not found in users", 404, err
	}

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", users.UidUser).First(&sessions).Error; err != nil {
		return "User ID not found in session", 404, err
	}
	sessions.DeviceId = session.DeviceID
	if sessions.ChangeDevice == nil {
		now := time.Now()
		sessions.ChangeDevice = &now
	}
	if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
		return "Failed update session", 500, err
	}
	if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
		return "Failed update session", 500, err
	}
	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ?", session.Phone).First(&trylimit).Error; err != nil {
		return "Phone not found in attempt", 404, err
	}
	trylimit.OtpAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return "Failed update attempt", 500, err
	}

	// Update User
	if users.Fullname == "" && users.Pin == "" {
		users.Status = "active"
		if err := r.Db.WithContext(ctx).Save(&users).Error; err != nil {
			return "Failed update status user", 500, err
		}
		return "User Register", 205, nil
	}
	return "Success Reset Device ID", 200, nil
}

func (r *session) LogoutSession(ctx context.Context, phone string, device *dto.DeviceId) (string, int, error) {
	var sessions model.Session
	var users model.User
	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", phone).First(&users).Error; err != nil {
		return "Phone not found in users", 404, err
	}

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ? AND device_id = ?", users.UidUser, device.DeviceId).First(&sessions).Error; err != nil {
		return "User ID And Device ID not found in session", 404, err
	}

	if sessions.LoggedOutAt == nil {
		now := time.Now()
		sessions.LoggedOutAt = &now
		sessions.TotalDevice = sessions.TotalDevice - 1
	}
	sessions.Status = false

	// Pisahkan perangkat yang ingin dihapus dari string "Device"
	devices := strings.Split(sessions.DeviceId, ",")
	var newDevices []string
	for _, d := range devices {
		if d != device.DeviceId {
			newDevices = append(newDevices, d)
		}
	}
	sessions.DeviceId = strings.Join(newDevices, ",")

	if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
		return "Failed update session", 500, err
	}

	return "Success Remove Session", 201, nil
}

func (r *session) CheckSession(ctx context.Context, uid_users string, device_id string, phone string, status int, msg string) (string, int, string, error) {
	var sessions model.Session
	otp := helper.GeneratePin(6)
	totalDevice := util.Getenv("TOTAL_DEVICE", "")
	intDevice, _ := strconv.ParseInt(totalDevice, 10, 16)
	int16Value := int16(intDevice)

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		return msg, 201, otp, nil
	}

	// Check all device id
	devices := strings.Split(sessions.DeviceId, ",")
	var isDevice bool
	for _, d := range devices {
		if d == device_id {
			isDevice = true
			break
		}
	}

	if isDevice && sessions.Status == true && sessions.LoggedOutAt == nil && status == 200 && sessions.TotalDevice <= int16Value {
		// User sudah ada device id yang sama ketika login
		return msg, status, otp, nil
	} else {
		// sessions.TotalDevice >= int16Value
		// User sudah melebihi 2 account akan kena limit dan already device
		return "Your Device Already 2 Account Login", 403, "Error", nil
	}
}
