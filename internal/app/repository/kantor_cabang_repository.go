package repository

import (
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KantorCabangRepository interface {
	Create(kantor *model.KantorCabang) error
	FindAll() ([]model.KantorCabang, error)
	ExistsById(c *fiber.Ctx, id uint) (bool, error)
	FindById(c *fiber.Ctx, id uint) (*model.KantorCabang, error)
}

type kantorCabangRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewKantorCabangRepository(db *gorm.DB, log *zap.Logger) KantorCabangRepository {
	return &kantorCabangRepository{db, log}
}

func (r *kantorCabangRepository) Create(kantor *model.KantorCabang) error {
	return r.db.Create(kantor).Error
}

func (r *kantorCabangRepository) ExistsById(c *fiber.Ctx, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(c.Context()).
		Model(&model.KantorCabang{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}
func (r *kantorCabangRepository) FindById(c *fiber.Ctx, id uint) (*model.KantorCabang, error) {
	var kantorCabang model.KantorCabang
	err := r.db.WithContext(c.Context()).
		Model(&model.KantorCabang{}).Where("id = ?", id).First(&kantorCabang).Error
	return &kantorCabang, err
}

func (r *kantorCabangRepository) FindAll() ([]model.KantorCabang, error) {
	var list []model.KantorCabang
	err := r.db.Find(&list).Error
	return list, err
}
