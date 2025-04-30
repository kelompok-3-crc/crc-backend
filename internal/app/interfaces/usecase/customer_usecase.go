package usecase

import (
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
)

type CustomerUsecase interface {
	// Login(c *fiber.Ctx, req dto.LoginRequest) (string, error)
	Create(c *fiber.Ctx, req dto.PredictionRequest) (*model.Customer, error)
}
