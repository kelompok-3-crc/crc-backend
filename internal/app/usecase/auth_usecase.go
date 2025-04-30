package usecase

import (
	"errors"
	"fmt"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/interfaces/repository"
	"ml-prediction/internal/app/interfaces/usecase"
	"ml-prediction/internal/app/model"
	"ml-prediction/pkg/jwt"
	"ml-prediction/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) usecase.AuthUsecase {
	return &authUsecase{userRepo}
}

func (s *authUsecase) Login(c *fiber.Ctx, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByNIP(req.NIP)
	if err != nil {

		return "", errors.New("kredensial login tidak valid")
	}

	err = utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return "", errors.New("password salah")
	}
	accessToken, err := jwt.GenerateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %w", err)
	}

	return accessToken, nil
}

func (s *authUsecase) CreateUser(c *fiber.Ctx, req dto.CreateRequest) (*model.User, error) {

	exists, err := s.userRepo.ExistsByNama(c, req.Nama)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user dengan nama yang diberikan telah ada")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Nama:     req.Nama,
		Role:     req.Role,
		Password: string(hashed),
	}

	user, err = s.userRepo.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
