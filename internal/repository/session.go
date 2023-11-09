package repository

import (
	"context"
	"fmt"
	"main/helper"
	dto "main/internal/dto"
	model "main/internal/model"
	"main/package/util"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Session interface {
	CreateSession(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, error)
	UpdateSession(c echo.Context, ctx context.Context, users model.User, session *dto.ReqSessionReset) (string, int, error)
	LogoutSession(c echo.Context, ctx context.Context, user model.User) (string, int, error)
	CheckSession(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, string, error)
	CheckSessionPin(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, error)
	CheckSessionReset(c echo.Context, ctx context.Context, uid_users string, phone *dto.CheckSession) (string, int, error)
}

type session struct {
	Db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *session {
	return &session{
		db,
	}
}

func (r *session) CreateSession(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, error) {
	var sessions model.Session

	totalDevice := util.Getenv("TOTAL_DEVICE", "")
	intDevice, _ := strconv.ParseInt(totalDevice, 10, 16)
	int16Value := int16(intDevice)

	// Header Application
	apps := c.Request().Header.Get("Application")
	device_id := c.Request().Header.Get("DeviceId")

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		// Create sesssion for new user initial
		newSession := model.Session{
			UserId:       uid_users,
			DeviceId:     device_id,
			LoginInAt:    time.Now(),
			TotalDevice:  1,
			Application:  apps,
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

	// Check all Apps id
	applications := strings.Split(sessions.Application, ",")
	var isApps bool
	for _, d := range applications {
		if d == apps {
			isApps = true
			break
		}
	}

	if isDevice == false && isApps == false && sessions.Status == true && sessions.LoggedOutAt == nil && status == 201 {
		// User device A sudah login, tetapi ada device B ingin login maka masih bisa
		sessions.TotalDevice = sessions.TotalDevice + 1
		if !strings.Contains(sessions.DeviceId, device_id) {
			sessions.DeviceId = sessions.DeviceId + "," + device_id
		}
		if !strings.Contains(sessions.Application, apps) {
			sessions.Application = sessions.Application + "," + apps
		}
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed create session", 500, err
		}
		return msg, status, nil
	} else if isDevice && isApps == false && sessions.Status == true && sessions.LoggedOutAt == nil && status == 201 {
		// Device sama, tetapi beda apps
		sessions.TotalDevice = sessions.TotalDevice + 1
		if !strings.Contains(sessions.Application, apps) {
			sessions.Application = sessions.Application + "," + apps
		}
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed create session", 500, err
		}
		return msg, status, nil
	} else if sessions.LoggedOutAt != nil && sessions.TotalDevice <= int16Value && sessions.Status == false {
		// User sudah logout di device a tetapi ingin login di device b
		sessions.Status = true
		sessions.DeviceId = device_id
		sessions.Application = apps
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

func (r *session) UpdateSession(c echo.Context, ctx context.Context, users model.User, session *dto.ReqSessionReset) (string, int, error) {
	var trylimit model.Attempt

	// Header Application
	apps := c.Request().Header.Get("Application")
	device_id := c.Request().Header.Get("DeviceId")

	var sessions model.Session
	if err := r.Db.WithContext(ctx).Where("user_id = ?", users.UidUser).Find(&sessions).Error; err != nil {
		return "Not Found Session to user id", 404, err
	}

	// Pisahkan nilai "device_id" dan "application" menjadi array
	deviceIDs := strings.Split(sessions.DeviceId, ",")
	applications := strings.Split(sessions.Application, ",")

	// Cari posisi name application dalam array "applications"
	var qgamesIndex = 0
	for i, app := range applications {
		if app == apps {
			qgamesIndex = i
			break
		}
	}

	if qgamesIndex != -1 {
		// Update nilai "device_id" yang sesuai dengan name application
		deviceIDs[qgamesIndex] = device_id
		// Gabungkan kembali array "device_id" menjadi string
		updatedDeviceID := strings.Join(deviceIDs, ",")
		// Update nilai "device_id" kembali ke database
		sessions.DeviceId = updatedDeviceID
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed Update Session", 500, err
		}
	}
	// Update trylimit
	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ?", session.Phone).First(&trylimit).Error; err != nil {
		return "Phone not found in attempt", 404, err
	}
	trylimit.OtpAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
		return "Failed update attempt", 500, err
	}

	// // Update User
	if users.Fullname == "" && users.Pin == "" {
		return "User Register", 205, nil
	}
	return "Success Reset Device ID", 200, nil
}

func (r *session) LogoutSession(c echo.Context, ctx context.Context, user model.User) (string, int, error) {
	var sessions model.Session
	totalDevice := util.Getenv("TOTAL_DEVICE", "")
	intDevice, _ := strconv.ParseInt(totalDevice, 10, 16)
	intValue := int(intDevice)

	// Header Application
	apps := c.Request().Header.Get("Application")
	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", user.UidUser).First(&sessions).Error; err != nil {
		return "User ID not found in session", 404, err
	}

	deviceIDs := strings.Split(sessions.DeviceId, ",")
	applications := strings.Split(sessions.Application, ",")

	var isApps bool
	var qgamesIndex = 0
	for i, app := range applications {
		if app == apps {
			isApps = true
			qgamesIndex = i
			break
		}
	}
	fmt.Println("apaa tuh if ? ", isApps)
	fmt.Println("apaa tuh if ? ", len(deviceIDs) == intValue)
	fmt.Println("apaa tuh if ? ", len(applications) == intValue)
	fmt.Println("--------------")
	fmt.Println("apaa tuh else if ? ", isApps)
	fmt.Println("apaa tuh else if ? ", len(deviceIDs) == 1)
	fmt.Println("apaa tuh else if ? ", len(applications) == intValue)
	fmt.Println("--------------")
	fmt.Println("apaa tuh else if 2 ? ", isApps)
	fmt.Println("apaa tuh else if 2 ? ", len(deviceIDs) == 1)
	fmt.Println("apaa tuh else if 2 ? ", len(applications) == 1)

	if isApps && len(deviceIDs) == intValue && len(applications) == intValue {
		fmt.Println("msk 1")
		deviceIDs[qgamesIndex] = ""
		applications[qgamesIndex] = ""
		updatedDeviceID := strings.Join(deviceIDs, "")
		updatedAppsID := strings.Join(applications, "")
		sessions.DeviceId = updatedDeviceID
		sessions.Application = updatedAppsID
		sessions.TotalDevice = sessions.TotalDevice - 1
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed Update Session", 500, err
		}
	} else if isApps && len(deviceIDs) == 1 && len(applications) == intValue {
		fmt.Println("msk 2")
		applications[qgamesIndex] = ""
		deviceIDs[qgamesIndex] = ""
		updatedAppsID := strings.Join(applications, "")
		sessions.Application = updatedAppsID
		sessions.TotalDevice = sessions.TotalDevice - 1
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed Update Session", 500, err
		}
	} else if isApps && len(deviceIDs) == 1 && len(applications) == 1 {
		fmt.Println("msk 3")
		applications[qgamesIndex] = ""
		deviceIDs[qgamesIndex] = ""
		updatedDeviceID := strings.Join(deviceIDs, "")
		updatedAppsID := strings.Join(applications, "")
		now := time.Now()
		sessions.DeviceId = updatedDeviceID
		sessions.Application = updatedAppsID
		sessions.LoggedOutAt = &now
		sessions.Status = false
		sessions.TotalDevice = sessions.TotalDevice - 1
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed Update Session", 500, err
		}
	} else {
		return "Not Found Aplication", 404, nil
	}

	return "Success Reset Device ID", 200, nil
}

func (r *session) CheckSession(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, string, error) {
	var sessions model.Session
	otp := helper.GeneratePin(6)
	totalDevice := util.Getenv("TOTAL_DEVICE", "")
	intDevice, _ := strconv.ParseInt(totalDevice, 10, 16)
	int16Value := int16(intDevice)
	var isDevice bool
	var isApps bool

	// Header
	apps := c.Request().Header.Get("Application")
	device_id := c.Request().Header.Get("DeviceId")

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		return msg, status, otp, nil
	}

	// Check all device id
	devices := strings.Split(sessions.DeviceId, ",")
	var positionDevice int
	for i, d := range devices {
		if d == device_id {
			isDevice = true
			positionDevice = i
			break
		}
	}

	// Check name application
	application := strings.Split(sessions.Application, ",")
	var positionApps int
	for i, d := range application {
		if d == apps {
			isApps = true
			positionApps = i
			break
		}
	}

	if isDevice && isApps && positionDevice == positionApps && sessions.Status == true && sessions.LoggedOutAt == nil && sessions.TotalDevice <= int16Value {
		// Device ID harus sesuai dengan application
		// User sudah ada device id yang sama ketika login
		return msg, status, otp, nil
	} else if sessions.TotalDevice >= int16Value {
		// User sudah melebihi 2 account akan kena limit dan already device
		return "Your Device Already 2 Account Login", 403, "Error", nil
	} else if isDevice == false && isApps == false && sessions.Status == true && sessions.LoggedOutAt == nil && sessions.TotalDevice <= int16Value {
		// Device A sudah ada, tetapi Device B ingin login maka wajib otp jika ingin login
		return msg, 201, otp, nil
	} else if isDevice && isApps == false && sessions.Status == true && sessions.LoggedOutAt == nil && sessions.TotalDevice <= int16Value {
		// Sama menggunakan Device A tetapi beda Apps, wajib OTP
		return msg, 201, otp, nil
	} else if sessions.Status == false && sessions.LoggedOutAt != nil && sessions.TotalDevice == 0 {
		// User logout semuanya
		return msg, 201, otp, nil
	} else {
		// default
		return "Device Already Login", 403, "Error", nil
	}
}

func (r *session) CheckSessionPin(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, error) {
	var sessions model.Session
	var isDevice bool
	var isApps bool

	// Header Application
	apps := c.Request().Header.Get("Application")
	device_id := c.Request().Header.Get("DeviceId")

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		return msg, status, nil
	}

	// Check all device id
	devices := strings.Split(sessions.DeviceId, ",")
	var positionDevice int
	for i, d := range devices {
		if d == device_id {
			isDevice = true
			positionDevice = i
			break
		}
	}

	// Check name application
	application := strings.Split(sessions.Application, ",")
	var positionApps int
	for i, d := range application {
		if d == apps {
			isApps = true
			positionApps = i
			break
		}
	}

	if isDevice && isApps && positionDevice == positionApps && sessions.Status == true && sessions.LoggedOutAt == nil {
		// User sudah ada device id yang sama ketika login
		return msg, status, nil
	} else if isDevice == false {
		// default
		return "Not Found Device ID", 403, nil
	} else {
		return "Not Found Application", 403, nil
	}
}

func (r *session) CheckSessionReset(c echo.Context, ctx context.Context, uid_users string, phone *dto.CheckSession) (string, int, error) {
	var sessions model.Session
	var isDevice bool
	var isApps bool

	// Header Application
	apps := c.Request().Header.Get("Application")
	device_id := c.Request().Header.Get("DeviceId")

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		return "User Not Found Session", 404, err
	}

	// Check all device id
	devices := strings.Split(sessions.DeviceId, ",")
	for _, d := range devices {
		if d == device_id {
			isDevice = true
			break
		}
	}

	// Check name application
	application := strings.Split(sessions.Application, ",")
	for _, d := range application {
		if d == apps {
			isApps = true
			break
		}
	}

	if phone.Reset == "device-id" && isDevice == false && isApps == true {
		// if device id tidak ada di session artinya baru, maka akan validasi otp
		return "Send OTP For New Device ID", 201, nil
	} else if phone.Reset == "pin" && isDevice == true && isApps == true {
		return "Send OTP For New PIN", 201, nil
	} else {
		return "Not Found Application or haven't new Device ID", 403, nil
	}
}
