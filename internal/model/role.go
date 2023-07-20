package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID      int16     `json:"id" gorm:"serial;primaryKey"`
	UidRole uuid.UUID `json:"uid_role" gorm:"type:char(36);not_null;unique"`
	Name    string    `json:"name" gorm:"varchar;not_null;unique"`
	Desc    string    `json:"desc" gorm:"varchar;not_null;unique"`
	Data    string    `json:"data" gorm:"varchar"`
	Status  string    `json:"status" gorm:"varchar"`
	Common
}

func (role *Role) BeforeSave(tx *gorm.DB) (err error) {
	if role.UidRole == uuid.Nil {
		role.UidRole = uuid.NewV4()
	}
	return nil
}
