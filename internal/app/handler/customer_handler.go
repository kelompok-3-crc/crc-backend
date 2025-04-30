package handler

import (
	"log"
	"ml-prediction/config"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/interfaces/usecase"

	"ml-prediction/pkg/response"
	"ml-prediction/pkg/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	CustomerUsecase usecase.CustomerUsecase
	cfg             config.Configuration
	val             *validator.Validate
}

func NewCustomerHandler(CustomerService usecase.CustomerUsecase, cfg config.Configuration, val *validator.Validate) *CustomerHandler {
	return &CustomerHandler{CustomerService, cfg, val}
}

func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req dto.PredictionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Input tidak sesuai", err.Error())
	}
	log.Println(req)

	if err := h.val.Struct(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errors := validation.MapValidationErrors(errs)
			return response.ErrorValidation(c, fiber.StatusBadRequest, "Kesalahan Validasi", errors)
		}
		return response.Error(c, fiber.StatusBadRequest, "Kesalahan Validasi", err.Error())
	}

	user, err := h.CustomerUsecase.Create(c, req)
	if err != nil {
		return response.Error(c, fiber.StatusConflict, "Gagal registrasi", err.Error())
	}

	return response.Success(c, "Customer berhasil dibuat!!", user)
}
