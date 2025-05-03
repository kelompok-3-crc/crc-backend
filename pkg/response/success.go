package response

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessCreated(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// SuccessResponseModel represents a standard success response for Swagger
type SuccessResponseModel struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponseModel represents an error response for Swagger
type ErrorResponseModel struct {
	Success bool        `json:"success" example:"false"`
	Error   string      `json:"error" example:"Error occurred"`
	Message interface{} `json:"message,omitempty"`
}

// LoginResponseModel represents a login success response for Swagger
type LoginResponseModel struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login berhasil!"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
