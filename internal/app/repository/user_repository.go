package repository

import (
	"fmt"
	"ml-prediction/internal/app/model"
	"ml-prediction/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByNIP(nip string) (*model.User, error)
	CreateUser(c *fiber.Ctx, user *model.User) (*model.User, error)
	ExistsByNama(c *fiber.Ctx, nama string) (bool, error)
	FindByNIPWithTx(tx *gorm.DB, nip string) (*model.User, error)
}
type userRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserRepo(db *gorm.DB, log *zap.Logger) UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (r *userRepository) FindByNIP(nip string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("KantorCabang").Where("nip = ?", nip).First(&user).Error
	return &user, err
}

func (r *userRepository) CreateUser(c *fiber.Ctx, user *model.User) (*model.User, error) {

	for {
		nip := utils.GenerateRandomNIP(10)
		var count int64
		err := r.db.WithContext(c.Context()).Model(&model.User{}).Where("nip = ?", nip).Count(&count).Error
		if err != nil {
			return nil, err
		}
		if count == 0 {
			user.NIP = nip
			break
		}
	}
	tx := r.db.WithContext(c.Context())

	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}

	var createdUser model.User
	if err := tx.Preload("KantorCabang").Where("id = ?", user.ID).First(&createdUser).Error; err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (r *userRepository) ExistsByNama(c *fiber.Ctx, nama string) (bool, error) {
	var count int64
	err := r.db.WithContext(c.Context()).
		Model(&model.User{}).
		Where("nama = ?", nama).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) FindByNIPWithTx(tx *gorm.DB, nip string) (*model.User, error) {
	var user model.User
	err := tx.Preload("KantorCabang").Where("nip = ?", nip).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user dengan NIP %s tidak ditemukan", nip)
		}
		return nil, err
	}
	return &user, nil
}
