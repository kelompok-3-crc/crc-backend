package handler

import (
	"ml-prediction/config"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/usecase"
	"ml-prediction/pkg/helper"
	"ml-prediction/pkg/response"
	"ml-prediction/pkg/validation"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TargetHandler interface {
	GetTargetSummary(c *fiber.Ctx) error
	CreateTargetTahunan(c *fiber.Ctx) error
}
type targetHandler struct {
	Usecase usecase.TargetUsecase
	cfg     config.Configuration
	val     *validator.Validate
}

func NewTargetHandler(u usecase.TargetUsecase,
	cfg config.Configuration,
	val *validator.Validate) TargetHandler {
	return &targetHandler{Usecase: u, cfg: cfg, val: val}
}

func (h *targetHandler) CreateTargetTahunan(c *fiber.Ctx) error {
	var req dto.TargetProdukTahunanRequest

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

	nip := c.Locals("nip").(string)

	if err := h.Usecase.CreateOrUpdateTargetTahunan(nip, &req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Gagal membuat target tahunan", err.Error())
	}

	return response.Success(c, "Target tahunan berhasil dibuat/diperbarui", nil)
}

func (h *targetHandler) GetTargetSummary(c *fiber.Ctx) error {
	// Extract user ID and role from context (set by auth middleware)
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", "User not authenticated")
	}
	userNIP := c.Locals("nip").(string)
	if userNIP == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", "User not authenticated")
	}

	var userID uint
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case string:
		idInt, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Invalid user ID format", err.Error())
		}
		userID = uint(idInt)
	default:
		return response.Error(c, fiber.StatusInternalServerError, "Invalid user ID format", "Unexpected type")
	}

	currentTime := time.Now()

	monthStr := c.Query("month")
	month := currentTime.Month()
	if monthStr != "" {
		monthInt, err := strconv.Atoi(monthStr)
		if err != nil || monthInt < 1 || monthInt > 12 {
			return response.Error(c, fiber.StatusBadRequest, "Invalid month", "Month must be between 1 and 12")
		}
		month = time.Month(monthInt)
	}

	yearStr := c.Query("year")
	year := currentTime.Year()
	if yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid year", "Year must be a valid number")
		}
	}

	// Get target summary
	summary, err := h.Usecase.GetTargetSummary(c.Context(), userID, userNIP, int(month), year)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get target summary", err.Error())
	}

	return response.Success(c, "Mengambil data profil berhasil", summary)
}
