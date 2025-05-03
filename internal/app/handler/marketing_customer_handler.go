package handler

import (
	"fmt"
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

type MarketingCustomerHandler interface {
	UpdateCustomerStatus(c *fiber.Ctx) error
	GetMonthlyMonitoring(c *fiber.Ctx) error
	GetMonthlyMonitoringMarketing(c *fiber.Ctx) error
	GetProductPerformance(c *fiber.Ctx) error
}
type marketingCustomerHandler struct {
	usecase usecase.MarketingCustomerUsecase
	cfg     config.Configuration
	val     *validator.Validate
}

func NewMarketingCustomerHandler(
	usecase usecase.MarketingCustomerUsecase,
	cfg config.Configuration,
	val *validator.Validate,
) MarketingCustomerHandler {
	return &marketingCustomerHandler{
		usecase: usecase,
		cfg:     cfg,
		val:     val,
	}
}

func (h *marketingCustomerHandler) UpdateCustomerStatus(c *fiber.Ctx) error {
	cif := c.Params("cif")
	if cif == "" {
		return response.Error(c, fiber.StatusBadRequest, "CIF tidak valid", "CIF harus diisi")
	}

	var req dto.UpdateCustomerStatusRequest
	if err := c.BodyParser(&req); err != nil {
		errors := helper.MapUnmarshalErrors(err)
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Format JSON tidak valid", errors)
	}
	req.CIF = cif

	if err := h.val.Struct(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errors := validation.MapValidationErrors(errs, &req)
			return response.ErrorValidation(c, fiber.StatusBadRequest, "Kesalahan Validasi", errors)
		}
		return response.Error(c, fiber.StatusBadRequest, "Kesalahan Validasi", err.Error())
	}

	marketingNIP := c.Locals("nip").(string)
	if err := h.usecase.UpdateCustomerStatus(&req, marketingNIP); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Gagal mengassign customer", err.Error())
	}

	return response.Success(c, "Customer berhasil diassign", nil)
}

func (h *marketingCustomerHandler) GetMonthlyMonitoring(c *fiber.Ctx) error {
	var req dto.MonitoringRequest
	if err := c.QueryParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format request tidak valid", err.Error())
	}

	if err := h.val.Struct(&req); err != nil {
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Validasi gagal", validation.MapValidationErrors(err, &req))
	}

	result, err := h.usecase.GetMonthlyMonitoring(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan data monitoring", err.Error())
	}

	return response.Success(c, "Berhasil mendapatkan data monitoring", result)
}
func (h *marketingCustomerHandler) GetMonthlyMonitoringMarketing(c *fiber.Ctx) error {

	nip := c.Locals("nip").(string)
	if nip == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Tidak terotorisasi", "Autentikasi gagal")
	}

	var req dto.MonitoringRequest
	if err := c.QueryParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format request tidak valid", err.Error())
	}

	req.NIP = nip

	if err := h.val.Struct(&req); err != nil {
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Validasi gagal", validation.MapValidationErrors(err, &req))
	}

	result, err := h.usecase.GetMonthlyMonitoringMarketing(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan data target", err.Error())
	}

	return response.Success(c, "Berhasil mendapatkan data target", result)
}

func (h *marketingCustomerHandler) GetProductPerformance(c *fiber.Ctx) error {
	currentYear := time.Now().Year()
	req := dto.ProductPerformanceRequest{
		StartDate: fmt.Sprintf("%d-01-01", currentYear),
		EndDate:   fmt.Sprintf("%d-12-31", currentYear),
		GroupBy:   "month",
	}

	// Explicitly extract each query parameter
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		productID, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid product_id", err.Error())
		}
		req.ProductID = uint(productID)
	}

	// Extract start_date if provided
	if startDate := c.Query("start_date"); startDate != "" {
		req.StartDate = startDate
	}

	// Extract end_date if provided
	if endDate := c.Query("end_date"); endDate != "" {
		req.EndDate = endDate
	}

	// Extract group_by if provided
	if groupBy := c.Query("group_by"); groupBy != "" {
		req.GroupBy = groupBy
	}

	if err := h.val.Struct(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errors := validation.MapValidationErrors(errs, &req)
			return response.ErrorValidation(c, fiber.StatusBadRequest, "Kesalahan Validasi", errors)
		}

		return response.Error(c, fiber.StatusBadRequest, "Kesalahan Validasi", err.Error())
	}

	result, err := h.usecase.GetProductPerformance(c.Context(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get product performance", err.Error())
	}

	return response.Success(c, "Successfully retrieved product performance", result)
}
