package model

import (
	"time"
)

type Session struct {
	ID           int16     `json:"id" gorm:"serial;primaryKey"`
	UserId       string    `json:"user_id" gorm:"uuid;not null;unique"`
	DeviceId     string    `json:"device_id" gorm:"text;not null"`
	LoginInAt    time.Time `json:"login_in_at" gorm:"timestamp"`
	Status       bool      `json:"status"`
	TotalDevice  int16     `json:"total_device" gorm:"int16;not null"`
	Application  string    `json:"application" gorm:"varchar;not null"`
	ChangeDevice *time.Time
	LoggedOutAt  *time.Time
}
