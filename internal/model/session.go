package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID           int16          `json:"id" gorm:"serial;primaryKey"`
	UserId       string         `json:"user_id" gorm:"uuid;not null"`
	DeviceId     string         `json:"device_id" gorm:"text;not null"`
	LoginInAt    time.Time      `json:"login_in_at" gorm:"timestamp"`
	ChangeDevice time.Time      `json:"change_device" gorm:"timestamp"`
	LoggedOutAt  gorm.DeletedAt `gorm:"index"`
}
