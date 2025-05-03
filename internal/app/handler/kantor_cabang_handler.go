package handler

import (
	"ml-prediction/config"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/usecase"
	"ml-prediction/pkg/helper"
	"ml-prediction/pkg/response"
	"ml-prediction/pkg/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type KantorCabangHandler struct {
	usecase usecase.KantorCabangUsecase
	cfg     config.Configuration
	val     *validator.Validate
}

func NewKantorCabangHandler(uc usecase.KantorCabangUsecase, cfg config.Configuration, val *validator.Validate) *KantorCabangHandler {
	return &KantorCabangHandler{uc, cfg, val}
}

func (h *KantorCabangHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateKantorCabangRequest

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
	result, err := h.usecase.Create(req)
	if err != nil {
		return response.Error(c, fiber.StatusConflict, "Gagal membuat kantor cabang", err.Error())
	}
	return response.SuccessCreated(c, "Kantor cabang berhasil dibuat", result)
}

func (h *KantorCabangHandler) GetAll(c *fiber.Ctx) error {
	results, err := h.usecase.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan data kantor cabang", err.Error())
	}
	return response.Success(c, "Berhasil mendapatkan data kantor cabang", results)
}
