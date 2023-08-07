package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID          int16  `json:"id" gorm:"serial;primaryKey"`
	UserId      string `json:"user_id" gorm:"uuid;not null"`
	DeviceId    string `json:"device_id" gorm:"text;not null"`
	LoginInAt   time.Time
	LoggedOutAt gorm.DeletedAt `gorm:"index"`
}
