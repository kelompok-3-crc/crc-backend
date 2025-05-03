package dto

import (
	"ml-prediction/internal/app/model"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CustomerSearchRequest struct {
	Page     int    `json:"page" query:"page"`
	Limit    int    `json:"limit" query:"limit"`
	Search   string `json:"search" query:"search"`
	SearchBy string `json:"search_by" query:"searchBy"`
}

type AssignedCustomerRequest struct {
	Page   int    `json:"page" query:"page"`
	Limit  int    `json:"limit" query:"limit"`
	Status string `json:"status" query:"status"`
}

type CustomerProductResponse struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	Ikon     string `json:"ikon"`
	Prediksi string `json:"prediksi"`
	Order    int    `json:"order"`
}

type Customer struct {
	Id                 uint64         `gorm:"primaryKey" json:"id"`
	Nama               string         `gorm:"type:varchar(100);not null" json:"nama"`
	CIF                string         `gorm:"type:varchar(50);not null;unique" json:"cif"`
	NomorRekening      string         `gorm:"type:varchar(50);not null;unique" json:"nomor_rekening"`
	NamaPerusahaan     string         `gorm:"type:varchar(50);not null;"  json:"nama_perusahaan"`
	ProdukEksisting    pq.StringArray `gorm:"type:varchar[]"  json:"produk_eksisting"`
	AktivitasTransaksi string         `gorm:"type:varchar(100)"  json:"aktivitas_transaksi"`
	NomorHp            string         `gorm:"type:varchar(20)"  json:"nomor_hp"`
	Segmen             string         `gorm:"type:varchar(20)"  json:"segmen"`
	Address            string         `gorm:"type:text"  json:"alamat"`
	Job                string         `gorm:"type:varchar(100)"  json:"pekerjaan"`
	Penghasilan        int64          `gorm:"type:int"  json:"penghasilan"`
	Umur               int            `gorm:"type:int"  json:"umur"`
	Gender             string         `gorm:"type:varchar(10)"  json:"gender"`
	StatusPerkawinan   bool           `gorm:"type:boolean"  json:"status_perkawinan"`
	Payroll            bool           `gorm:"type:boolean"  json:"payroll"`

	Status              string    `gorm:"type:sting"  json:"status"`
	Notes               string    `json:"catatan"`
	MarketingCustomerID uint      `json:"marketing_customer_id,omitempty" gorm:"column:marketing_customer_id"`
	MarketingID         uint      `json:"marketing_id,omitempty" gorm:"column:marketing_id"`
	MCCreatedAt         time.Time `json:"mc_created_at,omitempty" gorm:"column:mc_created_at"`

	CreatedAt time.Time                 ` json:"created_at"`
	UpdatedAt time.Time                 ` json:"updated_at"`
	DeletedAt gorm.DeletedAt            `gorm:"index" json:"deleted_at"`
	Produk    []CustomerProductResponse `json:"produk,omitempty" gorm:"-"`

	ClosedProductID uint          `json:"closed_produk_id,omitempty" gorm:"column:product_id"`
	ClosedProduk    model.Product `json:"closed_produk,omitempty" gorm:"-"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int64 `json:"total_pages"`
}
