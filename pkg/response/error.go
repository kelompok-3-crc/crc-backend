package response

import "github.com/gofiber/fiber/v2"

// ErrorResponse represents an error response
func Error(c *fiber.Ctx, status int, err string, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   err,
		"message": message,
	})
}

// ErrorValidationResponse represents an error validation response
func ErrorValidation(c *fiber.Ctx, status int, err string, validationErrors map[string]string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": err,
		"errors":  validationErrors,
	})
}

// ErrorFieldResponse represents an error field response
func ErrorField(c *fiber.Ctx, status int, err string, validationErrors map[string]interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": err,
		"errors":  validationErrors,
	})
}
