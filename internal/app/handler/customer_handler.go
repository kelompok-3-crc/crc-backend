package handler

import (
	"ml-prediction/config"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/usecase"
	"ml-prediction/pkg/helper"
	"ml-prediction/pkg/response"
	"ml-prediction/pkg/validation"
	"strconv"

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
	user, err := h.CustomerUsecase.Create(c, req)
	if err != nil {
		return response.Error(c, fiber.StatusConflict, "Gagal menambahkan customer", err.Error())
	}

	return response.Success(c, "Customer berhasil dibuat!!", user)
}

func (h *CustomerHandler) GetNewCustomers(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	searchBy := c.Query("searchBy", "all")

	req := dto.CustomerSearchRequest{
		Page:     page,
		Limit:    limit,
		Search:   search,
		SearchBy: searchBy,
	}
	NIP := c.Locals("nip").(string)

	customers, meta, err := h.CustomerUsecase.GetNewCustomers(c.Context(), NIP, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get customers: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Customers retrieved successfully",
		"data":    customers,
		"meta":    meta,
	})
}

func (h *CustomerHandler) GetAssignedCustomers(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "all")

	req := dto.AssignedCustomerRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
	}
	NIP := c.Locals("nip").(string)

	customers, meta, err := h.CustomerUsecase.GetAssignedCustomers(c.Context(), NIP, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get assigned customers: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Assigned customers retrieved successfully",
		"data":    customers,
		"meta":    meta,
	})
}

func (h *CustomerHandler) GetCustomerDetail(c *fiber.Ctx) error {

	customerID := c.Params("cif")

	NIP := c.Locals("nip").(string)

	customer, err := h.CustomerUsecase.GetCustomerDetail(c.Context(), NIP, customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get customer details: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Customer details retrieved successfully",
		"data":    customer,
	})
}
