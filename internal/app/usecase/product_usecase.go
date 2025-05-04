package usecase

import (
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"
)

type ProductUsecase interface {
	GetAllProducts() ([]model.Product, error)
}

// productUsecaseImpl is the implementation of ProductUsecase.
type productUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase creates a new instance of ProductUsecase.
func NewProductUsecase(r repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: r,
	}
}

// GetAllProducts returns all products.
func (u *productUsecase) GetAllProducts() ([]model.Product, error) {
	return u.productRepo.GetAllProducts()
}
