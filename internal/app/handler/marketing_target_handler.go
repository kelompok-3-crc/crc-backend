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
	"github.com/jinzhu/now"
)

type MarketingTargetHandler struct {
	targetUsecase usecase.MarketingTargetUsecase
	cfg           config.Configuration
	val           *validator.Validate
}

func NewMarketingTargetHandler(usecase usecase.MarketingTargetUsecase, cfg config.Configuration, val *validator.Validate) *MarketingTargetHandler {
	return &MarketingTargetHandler{
		targetUsecase: usecase,
		cfg:           cfg,
		val:           val,
	}
}
func (h *MarketingTargetHandler) AssignMarketingTarget(c *fiber.Ctx) error {
	var req dto.AssignMarketingTargetRequest
	req.Tahun = now.BeginningOfMonth().Year()
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

	userNIP := c.Locals("nip").(string)

	if err := h.targetUsecase.AssignBulkMarketingTarget(&req, userNIP); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Gagal membuat target", err.Error())
	}

	// var failedTargets []map[string]interface{}
	// for _, target := range req.Target {
	// 	singleTarget := &dto.SingleMarketingTarget{
	// 		Tahun:        req.Tahun,
	// 		Bulan:        req.Bulan,
	// 		MarketingNIP: req.MarketingNIP,
	// 		ProductID:    target.ProductID,
	// 		Amount:       target.Amount,
	// 	}

	// 	if err := h.targetUsecase.AssignMarketingTarget(singleTarget, userNIP); err != nil {
	// 		failedTargets = append(failedTargets, map[string]interface{}{
	// 			"product_id": target.ProductID,
	// 			"error":      err.Error(),
	// 		})
	// 	}
	// }

	return response.Success(c, "Semua target berhasil ditambahkan", map[string]interface{}{
		"total_targets": len(req.Target),
	})
}

func (h *MarketingTargetHandler) GetMarketingTargets(c *fiber.Ctx) error {
	var req dto.MonitoringRequest
	if err := c.QueryParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format request tidak valid", err.Error())
	}
	req.NIP = ""
	if err := h.val.Struct(&req); err != nil {
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Validasi gagal", validation.MapValidationErrors(err, &req))
	}

	result, err := h.targetUsecase.GetMarketingTargets(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan data target", err.Error())
	}

	return response.Success(c, "Berhasil mendapatkan data target", result)
}

func (h *MarketingTargetHandler) GetMarketingTargetsDetail(c *fiber.Ctx) error {
	var req dto.MonitoringRequest
	if err := c.QueryParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format request tidak valid", err.Error())
	}

	marketingNIP := c.Params("nip")
	if marketingNIP == "" {
		return response.Error(c, fiber.StatusBadRequest, "NIP tidak valid", "NIP harus diisi")
	}

	req.NIP = marketingNIP
	if err := h.val.Struct(&req); err != nil {
		return response.ErrorValidation(c, fiber.StatusBadRequest, "Validasi gagal", validation.MapValidationErrors(err, &req))
	}

	result, err := h.targetUsecase.GetMarketingTargets(&req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan data target", err.Error())
	}

	return response.Success(c, "Berhasil mendapatkan data target", result)
}
