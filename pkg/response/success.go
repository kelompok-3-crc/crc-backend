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
