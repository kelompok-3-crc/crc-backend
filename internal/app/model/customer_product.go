package model

import "time"

type CustomerProduct struct {
	CustomerID uint64   `gorm:"type:bigint;not null;column:customer_id" json:"customer_id"`
	ProductID  uint     `gorm:"type:bigint;not null;column:product_id" json:"product_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	Produk     Product  `gorm:"foreignKey:ProductID" json:"produk"`
	Order      int      `gorm:"column:order;default:0" json:"order"`
	PlafonMin  *uint64  `gorm:"" json:"plafon_min"`
	TenorMin   *string  `gorm:"type:varchar(100);" json:"tenor_min"`
	PlafonMax  *string  `gorm:"type:varchar(100);" json:"plafon_max"`
	TenorMax   *string  `gorm:"type:varchar(100);" json:"tenor_max"`
	CreatedAt  time.Time
}
