package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"type:varchar(60);not null;unique" json:"name"`
	NIP       string         `gorm:"type:varchar(20);not null;unique;column:nip" json:"nip"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role      string         `gorm:"type:varchar(20);not null"` // "bm" or "marketing"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
