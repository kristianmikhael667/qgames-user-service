package model

import uuid "github.com/satori/go.uuid"

type Role struct {
	ID      int16     `json:"id" gorm:"serial"`
	UidRole uuid.UUID `json:"uid_role" gorm:"uuid"`
	Name    string    `json:"name" gorm:"varchar"`
	Desc    string    `json:"desc" gorm:"varchar"`
	Data    string    `json:"data" gorm:"varchar"`
	Status  string    `json:"status" gorm:"varchar"`
	Common
}
