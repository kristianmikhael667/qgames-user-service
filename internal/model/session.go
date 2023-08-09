package model

import (
	"time"
)

type Session struct {
	ID           int16     `json:"id" gorm:"serial;primaryKey"`
	UserId       string    `json:"user_id" gorm:"uuid;not null;unique"`
	DeviceId     string    `json:"device_id" gorm:"text;not null;unique"`
	LoginInAt    time.Time `json:"login_in_at" gorm:"timestamp"`
	Status       bool      `json:"status"`
	ChangeDevice *time.Time
	LoggedOutAt  *time.Time
}
