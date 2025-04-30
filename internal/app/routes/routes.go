package routes

import (
	"ml-prediction/config"
	"ml-prediction/internal/app/handler"
	"ml-prediction/internal/app/repository"
	"ml-prediction/internal/app/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Register(api fiber.Router, db *gorm.DB, cfg config.Configuration, log *zap.Logger, val *validator.Validate) {
	userRepo := repository.NewUserRepo(db, log)
	authService := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authService, cfg, val)

	auth := api.Group("/auth")

	auth.Post("/login", authHandler.Login)
	auth.Post("/create", authHandler.CreateUser)

	customerRepo := repository.NewCustomerRepo(db, log)
	customerService := usecase.NewcustomerUsecase(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService, cfg, val)

	predict := api.Group("/predictions")
	predict.Post("/", customerHandler.CreateCustomer)
}
