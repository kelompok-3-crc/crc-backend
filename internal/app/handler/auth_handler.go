package handler

import (
	"ml-prediction/config"
	"ml-prediction/internal/app/usecase"
	"ml-prediction/pkg/helper"
	"ml-prediction/pkg/response"
	"ml-prediction/pkg/validation"

	dto "ml-prediction/internal/app/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthUsecase usecase.AuthUsecase
	cfg         config.Configuration
	val         *validator.Validate
}

func NewAuthHandler(authService usecase.AuthUsecase, cfg config.Configuration, val *validator.Validate) *AuthHandler {
	return &AuthHandler{authService, cfg, val}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		errors := helper.MapUnmarshalErrors(err)
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Format JSON tidak valid", errors)
	}
	if err := h.val.Struct(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errors := validation.MapValidationErrors(errs, &req)
			return response.ErrorValidation(c, fiber.StatusBadRequest, "Kesalahan Validasi", errors)
		}

		return response.Error(c, fiber.StatusBadRequest, "Kesalahan Validasi", err.Error())
	}

	res, err := h.AuthUsecase.Login(c, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Kredensial salah", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil!",
		"token":   res,
	})
}

func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateRequest

	if err := c.BodyParser(&req); err != nil {
		errors := helper.MapUnmarshalErrors(err)
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Format JSON tidak valid", errors)
	}
	if err := h.val.Struct(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errors := validation.MapValidationErrors(errs, &req)
			return response.ErrorValidation(c, fiber.StatusBadRequest, "Kesalahan Validasi", errors)
		}

		return response.Error(c, fiber.StatusBadRequest, "Kesalahan Validasi", err.Error())
	}
	user, err := h.AuthUsecase.CreateUser(c, req)
	if err != nil {
		return response.Error(c, fiber.StatusConflict, "Gagal registrasi", err.Error())
	}

	return response.Success(c, "User berhasil dibuat!!", user)
}
