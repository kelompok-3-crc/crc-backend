package model

import (
	"time"
)

type MarketingTargetBulanan struct {
	ID                    uint  `gorm:"primaryKey" json:"id"`
	Tahun                 int   `gorm:"type:int;not null" json:"tahun"`
	Bulan                 int   `gorm:"type:int;not null" json:"bulan"`
	TargetAmount          int64 `gorm:"type:bigint;not null" json:"target_amount"`
	MarketingID           uint  `gorm:"not null" json:"marketing_id"`
	ProductID             uint  `gorm:"not null" json:"product_id"`
	KantorCabangID        uint  `gorm:"not null" json:"kantor_cabang_id"`
	TargetProdukBulananID uint  `gorm:"not null" json:"target_produk_bulanan_id"`

	Marketing           User                `gorm:"foreignKey:MarketingID" json:"marketing"`
	Product             Product             `gorm:"foreignKey:ProductID" json:"product"`
	KantorCabang        KantorCabang        `gorm:"foreignKey:KantorCabangID" json:"kantor_cabang"`
	TargetProdukBulanan TargetProdukBulanan `gorm:"foreignKey:TargetProdukBulananID" json:"target_produk_bulanan"`

	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (MarketingTargetBulanan) TableName() string {
	return "marketing_target_bulanan"
}
