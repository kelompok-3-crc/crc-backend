package model

import (
	"time"

	"gorm.io/gorm"
)

type CustomerStatus string

const (
	CustomerStatusNew      CustomerStatus = "new"
	CustomerStatusClosed   CustomerStatus = "closed"
	CustomerStatusRejected CustomerStatus = "rejected"
)

type MarketingCustomer struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	CustomerID  uint64 `gorm:"uniqueIndex:idx_customer_unique" json:"customer_id"`
	MarketingID uint   `gorm:"not null" json:"marketing_id"`
	Status      string `gorm:"type:varchar(20);default:'new'" json:"status"`
	ProductID   *uint  `gorm:"null" json:"product_id"`
	Amount      *int64 `gorm:"null" json:"amount"`
	Notes       string `gorm:"type:text" json:"notes"`

	Customer  Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	Marketing User     `gorm:"foreignKey:MarketingID" json:"marketing"`
	Product   *Product `gorm:"foreignKey:ProductID" json:"product"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
