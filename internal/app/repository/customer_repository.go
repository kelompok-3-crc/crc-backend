package repository

import (
	"ml-prediction/internal/app/interfaces/repository"
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type customerRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCustomerRepo(db *gorm.DB, log *zap.Logger) repository.CustomerRepository {
	return &customerRepository{
		db:  db,
		log: log,
	}
}

func (r *customerRepository) Create(c *fiber.Ctx, user *model.Customer) (*model.Customer, error) {
	error := r.db.WithContext(c.Context()).Create(user).Error
	return user, error
}

func (r *customerRepository) ExistsByCif(c *fiber.Ctx, cif string) (bool, error) {
	var count int64
	err := r.db.WithContext(c.Context()).
		Model(&model.Customer{}).
		Where("cif = ?", cif).
		Count(&count).Error
	return count > 0, err
}
