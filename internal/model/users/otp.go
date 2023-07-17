package model

import "time"

type Otp struct {
	ID        int16     `json:"id" gorm:"serial"`
	Phone     string    `json:"phone" gorm:"varchar"`
	Otp       string    `json:"otp" gorm:"varchar"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}
