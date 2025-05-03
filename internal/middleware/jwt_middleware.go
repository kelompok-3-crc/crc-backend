package middleware

import (
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))

func JWTMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Tidak memiliki izin",
				"message": "Header Authorization tidak ditemukan",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Tidak memiliki izin",
				"message": "Format otorisasi tidak valid",
			})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("metode penandatanganan yang tidak terduga")
			}
			return accessSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Tidak memiliki izin",
				"message": "Token tidak valid atau kedaluwarsa",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Tidak memiliki izin",
				"message": "Klaim token tidak valid",
			})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Dilarang",
				"message": "Peran tidak ditemukan dalam token",
			})
		}

		isAuthorized := false
		for _, r := range requiredRoles {
			if role == r {
				isAuthorized = true
				break
			}
		}
		if !isAuthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Dilarang",
				"message": "Anda tidak memiliki akses ke sumber daya ini",
			})
		}

		c.Locals("user_id", claims["sub"])
		c.Locals("role", role)
		c.Locals("nip", claims["nip"])

		return c.Next()
	}
}
