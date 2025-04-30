package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Id                 uint64         `gorm:"primaryKey"`
	Nama               string         `gorm:"type:varchar(100);not null"`
	CIF                string         `gorm:"type:varchar(50);not null;unique"`
	NamaPerusahaan     string         `gorm:"type:varchar(50);not null;"`
	ProdukEksisting    pq.StringArray `gorm:"type:varchar[]"`
	AktivitasTransaksi string         `gorm:"type:varchar(100)"`
	NomorHp            string         `gorm:"type:varchar(20)"`
	Segmen             string         `gorm:"type:varchar(20)"`
	Address            string         `gorm:"type:text"`
	Job                string         `gorm:"type:varchar(100)"`
	Penghasilan        int64          `gorm:"type:int"`
	Umur               int            `gorm:"type:int"`
	Gender             string         `gorm:"type:varchar(10)"`
	StatusPerkawinan   bool           `gorm:"type:boolean"`
	Payroll            bool           `gorm:"type:boolean"`
	TopProduk          pq.StringArray `gorm:"type:varchar[]"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
