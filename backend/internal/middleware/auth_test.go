package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret"

func testKeyfunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.ErrSignatureInvalid
	}
	return []byte(testSecret), nil
}

func generateTestJWT(secret string, userID string, expired bool) string {
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp.Unix(),
		"aud": "authenticated",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	token := generateTestJWT(testSecret, userID, false)

	app := fiber.New()
	app.Use(NewAuthMiddleware(testKeyfunc))
	app.Get("/test", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		claims := c.Locals("claims")
		return c.JSON(fiber.Map{"user_id": uid, "has_claims": claims != nil})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(NewAuthMiddleware(testKeyfunc))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	token := generateTestJWT(testSecret, "550e8400-e29b-41d4-a716-446655440000", true)

	app := fiber.New()
	app.Use(NewAuthMiddleware(testKeyfunc))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
