package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Customer struct {
	Id                 uint64             `gorm:"primaryKey" json:"id"`
	Nama               string             `gorm:"type:varchar(100);not null" json:"nama"`
	CIF                string             `gorm:"type:varchar(50);not null;unique" json:"cif"`
	NomorRekening      string             `gorm:"type:varchar(50);not null;unique" json:"nomor_rekening"`
	NamaPerusahaan     string             `gorm:"type:varchar(50);not null;"  json:"nama_perusahaan"`
	ProdukEksisting    pq.StringArray     `gorm:"type:varchar[]"  json:"produk_eksisting"`
	AktivitasTransaksi string             `gorm:"type:varchar(100)"  json:"aktivitas_transaksi"`
	NomorHp            string             `gorm:"type:varchar(20)"  json:"nomor_hp"`
	Segmen             string             `gorm:"type:varchar(20)"  json:"segmen"`
	Address            string             `gorm:"type:text"  json:"alamat"`
	Job                string             `gorm:"type:varchar(100)"  json:"pekerjaan"`
	Email              string             `gorm:"type:varchar(50)"  json:"email"`
	Penghasilan        int64              `gorm:"type:int"  json:"penghasilan"`
	Umur               int                `gorm:"type:int"  json:"umur"`
	Gender             string             `gorm:"type:varchar(10)"  json:"gender"`
	StatusPerkawinan   bool               `gorm:"type:boolean"  json:"status_perkawinan"`
	Payroll            bool               `gorm:"type:boolean"  json:"payroll"`
	CustomerProduk     []*CustomerProduct `gorm:"foreignKey:CustomerID" json:"customer_produk"`
	CreatedAt          time.Time          ` json:"created_at"`
	UpdatedAt          time.Time          ` json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
}
