package repository

import (
	"context"
	"main/helper"
	model "main/internal/model"
	pkgutil "main/package/util"
	"main/package/util/response"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Otp interface {
	SendOtp(ctx context.Context, phone string, sc int, otp string, trylimit model.Attempt, msg string) (string, int, error)
}

type otp struct {
	Db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) *otp {
	return &otp{
		db,
	}
}

func (r *otp) SendOtp(ctx context.Context, phone string, sc int, otp string, trylimit model.Attempt, msg string) (string, int, error) {

	hashedOtp, err := pkgutil.HashPassword(otp)
	if err != nil {
		return "Error Hash OTP", 500, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	if sc == 201 {
		// Sementara karena blm ada otp dari vendor maka kita buat send ke logs dlu
		log.Print("Success Send Otp with Number: "+phone+" Your OTP : "+otp, 201)

		newOtp := model.Otp{
			Phone:     phone,
			Otp:       hashedOtp,
			ExpiredAt: time.Now().Add(2 * time.Minute),
		}
		if err := r.Db.WithContext(ctx).Save(&newOtp).Error; err != nil {
			return "Failed create otp", 500, err
		}

		timenow := time.Now()
		curr := timenow.Format("2006-01-02")
		lastTest := trylimit.LastAttempt.Format("2006-01-02")

		if (trylimit.OtpAttempt >= 3) && (curr == lastTest) {
			log.Print("Otp already 3 times", 400)
			return "Otp already 3 times", 400, response.CustomErrorBuilder(400, "Error", "Otp already 3 times")
		} else if curr != lastTest {
			trylimit.OtpAttempt = 0
		}

		trylimit.OtpAttempt++
		trylimit.LastAttempt = time.Now()

		if err := r.Db.WithContext(ctx).Save(&trylimit).Error; err != nil {
			return "Failed update attempt", 500, err
		}
		msg_otp, sc_otp := helper.SendOtp(phone, otp)
		if sc_otp != 200 && sc_otp != 201 {
			log.Print("Error "+msg_otp, 400)
			return msg_otp, sc_otp, nil
		}
		log.Print("Error "+msg_otp+" unclomplate otp", 201)
	}
	return msg, sc, nil
}
