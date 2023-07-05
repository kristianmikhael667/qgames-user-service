package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Assign struct {
	ID          int16     `json:"id" gorm:"serial"`
	UidAssign   uuid.UUID `json:"uid_assign" gorm:"uuid"`
	Users       string    `json:"users" gorm:"uuid"`
	Roles       string    `json:"roles" gorm:"uuid"`
	Permissions string    `json:"permissions" gorm:"uuid"`
	Status      string    `json:"status" gorm:"varchar"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Common
}
