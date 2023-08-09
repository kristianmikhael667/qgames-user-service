package repository

import (
	"context"
	"main/helper"
	dto "main/internal/dto"
	model "main/internal/model"
	"time"

	"gorm.io/gorm"
)

type Session interface {
	CreateSession(ctx context.Context, uid_users string, device_id string, phone string, status int, msg string) (string, int, string, error)
	UpdateSession(ctx context.Context, sc int, msg string, session *dto.ReqSessionReset) (string, int, error)
	LogoutSession(ctx context.Context, phone string, device *dto.DeviceId) (string, int, error)
}

type session struct {
	Db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *session {
	return &session{
		db,
	}
}

func (r *session) CreateSession(ctx context.Context, uid_users string, device_id string, phone string, status int, msg string) (string, int, string, error) {
	var sessions model.Session
	otp := helper.GeneratePin(6)

	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ?", uid_users).First(&sessions).Error; err != nil {
		// Create sesssion for new user
		newSession := model.Session{
			UserId:    uid_users,
			DeviceId:  device_id,
			LoginInAt: time.Now(),
		}
		if err := r.Db.WithContext(ctx).Save(&newSession).Error; err != nil {
			return "Failed create session", 500, "Error", err
		}
		return msg, 201, otp, nil
	}
	// Already Device
	if err := r.Db.WithContext(ctx).Model(&model.Session{}).Where("device_id = ?", device_id).First(&sessions).Error; err != nil {
		return "Device Already Login", 403, "Error", nil
	}
	return msg, status, otp, nil
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
	sessions.ChangeDevice = time.Now()
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

	sessions.LoggedOutAt = time.Now()
	if err := r.Db.WithContext(ctx).Save(&sessions).Error; err != nil {
		return "Failed update session", 500, err
	}

	return "Success Remove Session", 201, nil
}
