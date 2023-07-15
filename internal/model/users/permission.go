package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Permission struct {
	ID            int16     `json:"id" gorm:"serial"`
	UidPermission uuid.UUID `json:"uid_permission" gorm:"uuid"`
	Name          string    `json:"name" gorm:"varchar"`
	Slug          string    `json:"slug" gorm:"varchat"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Common
}
