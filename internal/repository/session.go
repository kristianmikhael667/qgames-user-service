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

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Session interface {
	CreateSession(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, error)
	UpdateSession(c echo.Context, ctx context.Context, users model.User, session *dto.ReqSessionReset) (string, int, error)
	LogoutSession(c echo.Context, ctx context.Context, phone string) (string, int, error)
	CheckSession(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, string, error)
	CheckSessionPin(c echo.Context, ctx context.Context, uid_users string, phone string, status int, msg string) (string, int, error)
	CheckSessionReset(c echo.Context, ctx context.Context, uid_users, phone string) (string, int, error)
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

func (r *session) LogoutSession(c echo.Context, ctx context.Context, phone string) (string, int, error) {
	var sessions model.Session
	var users model.User
	totalDevice := util.Getenv("TOTAL_DEVICE", "")
	intDevice, _ := strconv.ParseInt(totalDevice, 10, 16)
	intValue := int(intDevice)
	// Header Application
	apps := c.Request().Header.Get("Application")
	device_id := c.Request().Header.Get("DeviceId")

	if err := r.Db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", phone).First(&users).Error; err != nil {
		return "Phone not found in users", 404, err
	}

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", users.UidUser).First(&sessions).Error; err != nil {
		return "User ID not found in session", 404, err
	}

	// Device ID
	deviceId := device_id
	deviceIDSlice := strings.Split(deviceId, ",")
	var foundDeviceID string

	for _, device_id := range deviceIDSlice {
		if device_id == device_id {
			foundDeviceID = device_id
			break
		}
	}
	if foundDeviceID == "" {
		return "Device ID not found in session", 404, nil
	}

	// Application
	appId := sessions.Application
	appsSlice := strings.Split(appId, ",")
	var foundApps string

	for _, app := range appsSlice {
		if app == apps {
			foundApps = apps
			break
		}
	}
	if foundApps == "" {
		return "Application not found in session", 404, nil
	}

	if len(deviceIDSlice) >= intValue && len(appsSlice) >= intValue {
		if foundDeviceID != "" && foundApps != "" {
			// Device ID
			var newDeviceIDs []string
			for _, device_id := range deviceIDSlice {
				if device_id != device_id {
					newDeviceIDs = append(newDeviceIDs, device_id)
				}
			}
			updatedDeviceString := strings.Join(newDeviceIDs, ",")

			// Apps
			var newAppsIDs []string
			for _, app := range appsSlice {
				if app != apps {
					newAppsIDs = append(newAppsIDs, app)
				}
			}
			updatedAppsString := strings.Join(newAppsIDs, ",")

			// Update data
			sessions.DeviceId = updatedDeviceString
			sessions.Application = updatedAppsString
			sessions.TotalDevice = sessions.TotalDevice - 1
		}

		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed update session", 500, err
		}

		return "Success Remove Session Device 2", 201, nil
	} else if len(deviceIDSlice) == 1 && len(appsSlice) >= intValue {
		// Apps
		var newAppsIDs []string
		for _, app := range appsSlice {
			if app != apps {
				newAppsIDs = append(newAppsIDs, app)
			}
		}
		updatedAppsString := strings.Join(newAppsIDs, ",")

		// Update data
		sessions.Application = updatedAppsString
		sessions.TotalDevice = sessions.TotalDevice - 1

		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed update session", 500, err
		}
		return "Success Remove Session Apps 2", 201, nil
	} else {
		if sessions.LoggedOutAt == nil {
			now := time.Now()
			sessions.LoggedOutAt = &now
			sessions.TotalDevice = sessions.TotalDevice - 1
			sessions.DeviceId = ""
			sessions.Application = ""
		}
		sessions.Status = false
		if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
			return "Failed update session", 500, err
		}

		return "Success Remove Session Device 1", 201, nil
	}
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

	if isDevice && isApps && sessions.Status == true && sessions.LoggedOutAt == nil && sessions.TotalDevice <= int16Value {
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

	if isDevice && isApps && sessions.Status == true && sessions.LoggedOutAt == nil {
		// User sudah ada device id yang sama ketika login
		return msg, status, nil
	} else if isDevice == false {
		// default
		return "Not Found Device ID", 403, nil
	} else {
		return "Not Found Application", 403, nil
	}
}

func (r *session) CheckSessionReset(c echo.Context, ctx context.Context, uid_users, phone string) (string, int, error) {
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

	if isDevice == false && isApps == true {
		// if device id tidak ada di session artinya baru, maka akan validasi otp
		return "Send OTP For New Device ID", 201, nil
	} else {
		return "Not Found Application or haven't new Device ID", 403, nil
	}
}
