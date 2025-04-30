package repository

import (
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
)

type CustomerRepository interface {
	// FindByCif(cif string) (*model.Customer, error)
	Create(c *fiber.Ctx, user *model.Customer) (*model.Customer, error)
	ExistsByCif(c *fiber.Ctx, cif string) (bool, error)
}
