package usecase

import (
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"
)

type ProductUsecase interface {
	GetAllProducts() ([]model.Product, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(r repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: r,
	}
}

func (u *productUsecase) GetAllProducts() ([]model.Product, error) {
	return u.productRepo.GetAllProducts()
}
