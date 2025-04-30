package repository

import (
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	FindByNIP(nip string) (*model.User, error)
	CreateUser(c *fiber.Ctx, user *model.User) (*model.User, error)
	ExistsByNama(c *fiber.Ctx, nama string) (bool, error)
}
