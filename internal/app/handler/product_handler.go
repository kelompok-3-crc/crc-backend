package handler

import (
	"ml-prediction/internal/app/usecase"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles requests related to products.
type ProductHandler struct {
	ProductUsecase usecase.ProductUsecase
}

// NewProductHandler returns a new ProductHandler.
func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		ProductUsecase: u,
	}
}

// GetAllProducts handles GET requests to retrieve all products.
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.ProductUsecase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve products"})
	}
	return c.JSON(products)
}
