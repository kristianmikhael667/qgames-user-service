package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID            int16     `json:"id" gorm:"serial;primaryKey"`
	UidPermission uuid.UUID `json:"uid_permission" gorm:"type:char(36);not_null;unique"`
	Name          string    `json:"name" gorm:"varchar;not_null;unique"`
	Slug          string    `json:"slug" gorm:"varchar;not_null;unique"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Common
}

func (permission *Permission) BeforeSave(tx *gorm.DB) (err error) {
	if permission.UidPermission == uuid.Nil {
		permission.UidPermission = uuid.NewV4()
	}
	return nil
}
