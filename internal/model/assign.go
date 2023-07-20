package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Assign struct {
	ID          int16     `json:"id" gorm:"serial;primaryKey"`
	UidAssign   uuid.UUID `json:"uid_assign" gorm:"uuid"`
	Users       string    `json:"users" gorm:"uuid"`
	Roles       string    `json:"roles" gorm:"varchar"`
	Permissions string    `json:"permissions" gorm:"varchar"`
	Status      string    `json:"status" gorm:"varchar"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Common
}

func (role *Assign) BeforeSave(tx *gorm.DB) (err error) {
	if role.UidAssign == uuid.Nil {
		role.UidAssign = uuid.NewV4()
	}
	return nil
}
