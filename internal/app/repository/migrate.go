package repository

import (
	"log"

	"gorm.io/gorm"

	"ml-prediction/internal/app/model"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.Customer{})
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}
}
