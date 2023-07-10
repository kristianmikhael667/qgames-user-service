package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        int16     `json:"id" gorm:"serial"`
	UidUser   uuid.UUID `gorm:"type:char(36);primaryKey" json:"uid_user"`
	Fullname  string    `json:"fullname" gorm:"varchar;not_null"`
	Phone     string    `json:"phone" gorm:"varchar;not_null;unique"`
	Email     string    `json:"email" gorm:"varchar;not_null;unique"`
	Password  string    `json:"password" gorm:"varchar"`
	Pin       string    `json:"pin" gorm:"varchar"`
	Address   string    `json:"address" gorm:"text;not_null"`
	Profile   string    `json:"profile" gorm:"varchar"`
	Status    string    `json:"status" gorm:"varchar"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Common
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if user.UidUser == uuid.Nil {
		user.UidUser = uuid.NewV4()
	}
	return nil
}
