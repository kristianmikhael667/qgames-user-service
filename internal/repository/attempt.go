package repository

import (
	"context"
	"fmt"
	"main/internal/dto"
	model "main/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Attempt interface {
	CreateAttempt(ctx context.Context, phone string) (model.Attempt, int, string, error)
	UpdateAttemptOtp(ctx context.Context, phone string) (int, string, error)
	ResetAttemptOtp(ctx context.Context, payload *dto.RequestReset) (int, string, error)
	ResetAttemptPin(ctx context.Context, payload *dto.RequestReset) (int, string, error)
}

type attempt struct {
	Db *gorm.DB
}

func NewAttemptRepository(db *gorm.DB) *attempt {
	return &attempt{
		db,
	}
}

func (r *attempt) CreateAttempt(ctx context.Context, phone string) (model.Attempt, int, string, error) {
	var trylimit model.Attempt

	timenow := time.Now()

	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ? ", phone).First(&trylimit).Error; err != nil {
		// create try limit
		newAttemp := model.Attempt{
			Phone:       phone,
			PinAttempt:  0,
			OtpAttempt:  0,
			LastAttempt: timenow,
		}
		if err := r.Db.WithContext(ctx).Save(&newAttemp).Error; err != nil {
			return newAttemp, 500, "Failed create attemp", err
		}
		return newAttemp, 201, "Success Create Attempt for new users", nil
	}
	return trylimit, 200, "Attempt already", nil
}

func (r *attempt) UpdateAttemptOtp(ctx context.Context, phone string) (int, string, error) {
	var attemp model.Attempt

	phones := strings.Replace(phone, "+62", "0", -1)
	// phones = strings.Replace(phones, "62", "0", -1)

	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ?", phones).First(&attemp).Error; err != nil {
		return 404, "Error get phone attemp", err
	}
	attemp.OtpAttempt = 0
	if err := r.Db.WithContext(ctx).Save(&attemp).Error; err != nil {
		return 500, "Error update phone attemp", err
	}
	return 201, "Success update", nil
}

// Admin will reset otp
func (r *attempt) ResetAttemptOtp(ctx context.Context, payload *dto.RequestReset) (int, string, error) {
	var attemp model.Attempt

	phones := strings.Replace(payload.Phone, "+62", "0", -1)

	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ?", phones).First(&attemp).Error; err != nil {
		return 404, "Error get phone attemp", err
	}
	if attemp.OtpAttempt == 3 {
		attemp.OtpAttempt = 0
		if err := r.Db.WithContext(ctx).Save(&attemp).Error; err != nil {
			return 500, "Error update phone attemp", err
		}
		return 201, "Success Reset OTP", nil
	} else {
		return 403, "Your OTP have to up 3", nil
	}

}

// Admin will reset otp
func (r *attempt) ResetAttemptPin(ctx context.Context, payload *dto.RequestReset) (int, string, error) {
	var attemp model.Attempt

	phones := strings.Replace(payload.Phone, "+62", "0", -1)

	if err := r.Db.WithContext(ctx).Model(&model.Attempt{}).Where("phone = ?", phones).First(&attemp).Error; err != nil {
		return 404, "Error get phone attemp", err
	}
	fmt.Println("Adada ", attemp.PinAttempt)
	if attemp.PinAttempt == 3 {
		attemp.PinAttempt = 0
		if err := r.Db.WithContext(ctx).Save(&attemp).Error; err != nil {
			return 500, "Error update phone attemp", err
		}
		return 201, "Success Reset PIN", nil
	} else {
		return 403, "Your PIN have to up 3", nil
	}
}
