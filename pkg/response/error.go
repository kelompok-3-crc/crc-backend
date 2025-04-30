package response

import "github.com/gofiber/fiber/v2"

func Error(c *fiber.Ctx, status int, err string, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   err,
		"message": message,
	})
}

func ErrorValidation(c *fiber.Ctx, status int, err string, message map[string]string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   err,
		"message": message,
	})
}
