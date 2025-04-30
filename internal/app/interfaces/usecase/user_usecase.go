package usecase

import (
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
)

type AuthUsecase interface {
	Login(c *fiber.Ctx, req dto.LoginRequest) (string, error)
	CreateUser(c *fiber.Ctx, req dto.CreateRequest) (*model.User, error)
}
