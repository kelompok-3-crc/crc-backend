package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nama     string `gorm:"type:varchar(60);not null;unique" json:"name"`
	NIP      string `gorm:"type:varchar(20);not null;unique;column:nip" json:"nip"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(20);not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	KantorCabangID   *uint                    `gorm:"null" json:"kantor_cabang_id"`
	KantorCabang     KantorCabang             `gorm:"foreignKey:KantorCabangID" json:"kantor_cabang"`
	MarketingTargets []MarketingTargetBulanan `gorm:"foreignKey:MarketingID" json:"marketing_targets"`
	Products         []Product                `gorm:"many2many:customer_products;" json:"products"`
}
