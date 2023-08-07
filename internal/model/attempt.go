package model

import "time"

type Attempt struct {
	ID          int16     `json:"id" gorm:"serial"`
	Phone       string    `json:"phone" gorm:"varchar;not_null;unique"`
	PinAttempt  int64     `json:"pin_attempt" gorm:"int64"`
	OtpAttempt  int64     `json:"otp_attempt" gorm:"int64"`
	LastAttempt time.Time `json:"last_attempt"`
}
