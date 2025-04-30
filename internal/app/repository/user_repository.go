package repository

import (
	"ml-prediction/internal/app/interfaces/repository"
	"ml-prediction/internal/app/model"
	"ml-prediction/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserRepo(db *gorm.DB, log *zap.Logger) repository.UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (r *userRepository) FindByNIP(nip string) (*model.User, error) {
	var user model.User
	err := r.db.Where("nip = ?", nip).First(&user).Error
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

	error := r.db.WithContext(c.Context()).Create(user).Error
	return user, error
}

func (r *userRepository) ExistsByNama(c *fiber.Ctx, nama string) (bool, error) {
	var count int64
	err := r.db.WithContext(c.Context()).
		Model(&model.User{}).
		Where("nama = ?", nama).
		Count(&count).Error
	return count > 0, err
}
