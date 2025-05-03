package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"type:varchar(60);not null;unique" json:"nama"`
	Prediksi  string         `gorm:"type:varchar(50);" json:"prediksi"`
	Ikon      string         `gorm:"type:varchar(50);" json:"ikon"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// Customers []Customer     `gorm:"many2many:customer_products;" json:"customers"`
}
