package model

import (
	"time"

	"gorm.io/gorm"
)

type CustomerRecommendation struct {
	Id         uint64         `gorm:"primaryKey" json:"id"`
	PlafonMin  string         `gorm:"type:varchar(100);not null" json:"plafon_min"`
	TenorMin   string         `gorm:"type:varchar(100);not null" json:"tenor_min"`
	PlafonMax  string         `gorm:"type:varchar(100);not null" json:"plafon_max"`
	Tenormax   string         `gorm:"type:varchar(100);not null" json:"tenor_max"`
	CustomerId uint64         `gorm:"not null" json:"customer_id"`
	Customer   Customer       `gorm:"foreignKey:CustomerId" json:"customer"`
	ProductId  uint64         `gorm:"not null" json:"product_id"`
	Product    Product        `gorm:"foreignKey:ProductId" json:"product"`
	CreatedAt  time.Time      ` json:"created_at"`
	UpdatedAt  time.Time      ` json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
