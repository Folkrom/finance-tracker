package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// NewAdminMiddleware checks that the authenticated user has admin role
// in their JWT app_metadata. Must run after NewAuthMiddleware.
func NewAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(jwt.MapClaims)

		if !isAdmin(claims) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "admin access required",
			})
		}

		return c.Next()
	}
}

// IsAdmin checks if the current request is from an admin user.
func IsAdmin(c *fiber.Ctx) bool {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return false
	}
	return isAdmin(claims)
}

func isAdmin(claims jwt.MapClaims) bool {
	appMeta, ok := claims["app_metadata"].(map[string]any)
	if !ok {
		return false
	}
	role, ok := appMeta["role"].(string)
	if !ok {
		return false
	}
	return role == "admin"
}
