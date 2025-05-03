package model

import (
	"time"

	"gorm.io/gorm"
)

type KantorCabang struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Nama          string `gorm:"type:varchar(100);not null" json:"nama"`
	Users         []User
	TargetTahunan []TargetProdukTahunan `gorm:"foreignKey=KantorCabangID" json:"target_tahunan"`
	TargetBulanan []TargetProdukBulanan `gorm:"foreignKey=KantorCabangID" json:"target_bulanan"`

	CreatedAt time.Time      ` json:"created_at"`
	UpdatedAt time.Time      ` json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type TargetProdukTahunan struct {
	ID             uint  `gorm:"primaryKey" json:"id"`
	Tahun          int   `gorm:"not null" json:"tahun"`
	TargetAmount   int64 `gorm:"not null" json:"target_amount"`
	KantorCabangID uint  `gorm:"not null" json:"kantor_cabang_id"`
	ProductID      uint  `gorm:"not null" json:"product_id"`

	KantorCabang   KantorCabang          `gorm:"foreignKey=KantorCabangID" json:"kantor_cabang"`
	Product        Product               `gorm:"foreignKey=ProductID" json:"product"`
	BulananTargets []TargetProdukBulanan `gorm:"foreignKey:TargetTahunanID" json:"bulanan_targets"`

	CreatedAt time.Time      ` json:"created_at"`
	UpdatedAt time.Time      ` json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type TargetProdukBulanan struct {
	ID              uint  `gorm:"primaryKey" json:"id"`
	Tahun           int   `gorm:"type:int;not null" json:"tahun"`
	Bulan           int   `gorm:"type:int;not null" json:"bulan"`
	TargetAmount    int64 `gorm:"type:bigint;not null" json:"target_amount"`
	KantorCabangID  uint  `gorm:"not null" json:"kantor_cabang_id"`
	ProductID       uint  `gorm:"not null" json:"product_id"`
	TargetTahunanID uint  `gorm:"not null" json:"target_tahunan_id"`

	KantorCabang  KantorCabang        `gorm:"foreignKey:KantorCabangID" json:"kantor_cabang"`
	Product       Product             `gorm:"foreignKey:ProductID" json:"product"`
	TargetTahunan TargetProdukTahunan `gorm:"foreignKey:TargetTahunanID" json:"target_tahunan"`

	CreatedAt time.Time      ` json:"created_at"`
	UpdatedAt time.Time      ` json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Tabler interface {
	TableName() string
}

func (KantorCabang) TableName() string {
	return "kantor_cabang"
}
func (TargetProdukTahunan) TableName() string {
	return "target_produk_tahunan"
}
func (TargetProdukBulanan) TableName() string {
	return "target_produk_bulanan"
}
