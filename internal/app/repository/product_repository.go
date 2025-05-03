package repository

import (
	"fmt"
	"ml-prediction/internal/app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByPrediksi(name string) (*model.Product, error)
}
type productRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewProductRepo(db *gorm.DB, log *zap.Logger) ProductRepository {
	return &productRepository{
		db:  db,
		log: log,
	}
}
func (r *productRepository) FindByPrediksi(prediksi string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Where("prediksi = ?", prediksi).First(&product).Error; err != nil {
		return nil, fmt.Errorf("product with prediksi '%s' not found: %v", prediksi, err)
	}
	return &product, nil
}
